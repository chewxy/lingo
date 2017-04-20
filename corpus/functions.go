package corpus

import (
	"math"
	"strings"
	"unicode/utf8"

	"github.com/chewxy/lingo"
	"github.com/chewxy/lingo/treebank"
	"github.com/pkg/errors"
)

// GenerateCorpus creates a Corpus given a set of SentenceTag from a training set.
func GenerateCorpus(sentenceTags []treebank.SentenceTag) *Corpus {
	words := make([]string, 3)
	frequencies := make([]int, 3)

	words[0] = ""      // aka NULL, for when no word can be found
	frequencies[0] = 0 // no word is never found

	words[1] = "-UNKNOWN-"
	frequencies[1] = 0

	words[2] = "-ROOT-"
	frequencies[2] = 1

	knownWords := make(map[string]int)
	knownWords[""] = 0
	knownWords["-UNKNOWN-"] = 1
	knownWords["-ROOT-"] = 2

	maxWordLength := 0

	for _, sentenceTag := range sentenceTags {
		for _, lex := range sentenceTag.Sentence {
			id, ok := knownWords[lex.Value]
			if !ok {
				knownWords[lex.Value] = len(words)
				words = append(words, lex.Value)
				frequencies = append(frequencies, 1)

				runeCount := utf8.RuneCountInString(lex.Value)
				if runeCount > maxWordLength {
					maxWordLength = runeCount
				}
			} else {
				frequencies[id]++
			}
		}
	}

	var totals int
	for _, f := range frequencies {
		totals += f
	}

	return &Corpus{words, frequencies, knownWords, int64(len(words)), totals, maxWordLength}
}

// ViterbiSplit is a Viterbi algorithm for splitting words given a corpus
func ViterbiSplit(input string, c *Corpus) []string {
	s := strings.ToLower(input)
	probabilities := []float64{1.0}
	lasts := []int{0}

	runes := []int{}
	for i := range s {
		runes = append(runes, i)
	}
	runes = append(runes, len(s)+1)

	for i := range s {
		probs := make([]float64, 0)
		ls := make([]int, 0)

		// m := maxInt(0, i-c.maxWordLength)

		for j, r := range runes {
			if r > i {
				break
			}

			p, ok := c.WordProb(s[r : i+1])
			if !ok {
				// http://stackoverflow.com/questions/195010/how-can-i-split-multiple-joined-words#comment48879458_481773
				p = (math.Log(float64(1)/float64(c.totalFreq)) - float64(c.maxWordLength) - float64(1)) * float64(i-r) // note it should be i-r not j-i as per the SO post
			}
			prob := probabilities[j] * p

			probs = append(probs, prob)
			ls = append(ls, r)
		}

		maxProb := -math.SmallestNonzeroFloat64
		maxK := -1 << 63
		for j, p := range probs {
			if p > maxProb {
				maxProb = p
				maxK = ls[j]
			}
		}
		probabilities = append(probabilities, maxProb)
		lasts = append(lasts, maxK)
	}

	words := make([]string, 0)
	i := utf8.RuneCountInString(s)

	for i > 0 {
		start := lasts[i]
		words = append(words, s[start:i])
		i = start
	}

	// reverse it
	for i, j := 0, len(words)-1; i < j; i, j = i+1, j-1 {
		words[i], words[j] = words[j], words[i]
	}

	return words
}

// CosineSimilarity measures the cosine similarity of two strings.
func CosineSimilarity(a, b []string) float64 {
	countsA := make([]float64, 0)
	countsB := make([]float64, 0)
	uniques := make(map[string]int)

	// index the strings first
	for _, st := range a {
		s := strings.ToLower(st)
		id, ok := uniques[s]
		if !ok {
			uniques[s] = len(countsA)
			countsA = append(countsA, 1)
			countsB = append(countsB, 0) // create for countsB, but don't add
		} else {
			countsA[id]++
		}
	}

	for _, st := range b {
		s := strings.ToLower(st)
		id, ok := uniques[s]
		if !ok {
			uniques[s] = len(countsA)
			countsA = append(countsA, 0)
			countsB = append(countsB, 1)
		} else {
			countsB[id]++
		}
	}

	magA, err := mag(countsA)
	if err != nil {
		panic(err)
	}

	magB, err := mag(countsB)
	if err != nil {
		panic(err)
	}

	dotProd, err := dot(countsA, countsB)
	if err != nil {
		panic(err)
	}

	return dotProd / (magA * magB)

}

// DamerauLevenshtein calculates the Damerau-Levensthtein distance between two strings. See more at https://en.wikipedia.org/wiki/Damerau%E2%80%93Levenshtein_distance
func DamerauLevenshtein(s1 string, s2 string) (distance int) {
	// index by code point, not byte
	r1 := []rune(s1)
	r2 := []rune(s2)

	// the maximum possible distance
	inf := len(r1) + len(r2)

	// if one string is blank, we needs insertions
	// for all characters in the other one
	if len(r1) == 0 {
		return len(r2)
	}

	if len(r2) == 0 {
		return len(r1)
	}

	// construct the edit-tracking matrix
	matrix := make([][]int, len(r1))
	for i := range matrix {
		matrix[i] = make([]int, len(r2))
	}

	// seen characters
	seenRunes := make(map[rune]int)

	if r1[0] != r2[0] {
		matrix[0][0] = 1
	}

	seenRunes[r1[0]] = 0
	for i := 1; i < len(r1); i++ {
		deleteDist := matrix[i-1][0] + 1
		insertDist := (i+1)*1 + 1
		var matchDist int
		if r1[i] == r2[0] {
			matchDist = i
		} else {
			matchDist = i + 1
		}
		matrix[i][0] = minInt(minInt(deleteDist, insertDist), matchDist)
	}

	for j := 1; j < len(r2); j++ {
		deleteDist := (j + 1) * 2
		insertDist := matrix[0][j-1] + 1
		var matchDist int
		if r1[0] == r2[j] {
			matchDist = j
		} else {
			matchDist = j + 1
		}

		matrix[0][j] = minInt(minInt(deleteDist, insertDist), matchDist)
	}

	for i := 1; i < len(r1); i++ {
		var maxSrcMatchIndex int
		if r1[i] == r2[0] {
			maxSrcMatchIndex = 0
		} else {
			maxSrcMatchIndex = -1
		}

		for j := 1; j < len(r2); j++ {
			swapIndex, ok := seenRunes[r2[j]]
			jSwap := maxSrcMatchIndex
			deleteDist := matrix[i-1][j] + 1
			insertDist := matrix[i][j-1] + 1
			matchDist := matrix[i-1][j-1]
			if r1[i] != r2[j] {
				matchDist += 1
			} else {
				maxSrcMatchIndex = j
			}

			// for transpositions
			var swapDist int
			if ok && jSwap != -1 {
				iSwap := swapIndex
				var preSwapCost int
				if iSwap == 0 && jSwap == 0 {
					preSwapCost = 0
				} else {
					preSwapCost = matrix[maxInt(0, iSwap-1)][maxInt(0, jSwap-1)]
				}
				swapDist = i + j + preSwapCost - iSwap - jSwap - 1
			} else {
				swapDist = inf
			}
			matrix[i][j] = minInt(minInt(minInt(deleteDist, insertDist), matchDist), swapDist)
		}
		seenRunes[r1[i]] = i
	}

	return matrix[len(r1)-1][len(r2)-1]
}

// LongestCommonPrefix takes a slice of strings, and finds the longest common prefix
func LongestCommonPrefix(strs ...string) string {
	switch len(strs) {
	case 0:
		return "" // idiots
	case 1:
		return strs[0]
	}

	min := strs[0]
	max := strs[0]

	for _, s := range strs[1:] {
		switch {
		case s < min:
			min = s
		case s > max:
			max = s
		}
	}

	for i := 0; i < len(min) && i < len(max); i++ {
		if min[i] != max[i] {
			return min[:i]
		}
	}

	// In the case where lengths are not equal but all bytes
	// are equal, min is the answer ("foo" < "foobar").
	return min
}

/* The following two functions help in parsing a string into numbers. It's recommended you write abstractions over the functions*/

// StrsToInts converts a string slice into an int slice, with the help of NumberWords.
// The function assumes all helper words like "and" have been stripped.
// 		"One hundred and five" -> []string{"one", "hundred", "five"}
// This is a very primitive method, and doesn't take into account other words like "a hundred" or "a couple of hundred"
func StrsToInts(strs []string) (retVal []int, err error) {
	for _, s := range strs {
		intVal, ok := lingo.NumberWords[s]
		if !ok {
			return nil, errors.Errorf("Unable to parse the words %q as numbers", s)
		}

		if len(retVal) > 0 && intVal == 100 && retVal[len(retVal)-1] < 100 {
			retVal[len(retVal)-1] *= 100
		} else if len(retVal) > 0 && retVal[len(retVal)-1] < 1000 && intVal < 1000 {
			retVal[len(retVal)-1] += intVal
		} else {
			retVal = append(retVal, intVal)
		}
	}
	return
}

// CombineInts takes a int slice, and tries to make it one integer.
// It works by taking advantage of english - anything more than 1000 has a repeated pattern
// e.g.
// 		one hundred and fifty thousand two hundred and two
// there are 2 repeated patterns (one hundred and fifty) and  (two hundred and two)
//
// This allows us to repeatedly combine by addition or multiplication until there is one left
func CombineInts(ints []int) int {
	var total int
	for len(ints) > 0 {
		if len(ints) == 1 || ints[0] >= 1000 {
			last := ints[len(ints)-1]
			total += last
			ints = ints[0 : len(ints)-1] //pop it
		} else {
			if ints[1] < 1000 {
				// something went wrong
				panic("HELP!")
			}
			total += ints[0] * ints[1]
			ints = ints[2:]
		}
	}
	return total
}
