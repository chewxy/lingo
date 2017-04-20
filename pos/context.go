package pos

import (
	"strconv"

	"github.com/chewxy/lingo"
)

/*
A context is which word in the current state the POSTagger is in.
There are so far  5 contexts:
	- Previous previous word
	- previous word
	- current word
	- next word
	- next next word

For each context we have 8 features:
	- word (lower case)
	- lemma
	- cluster
	- shape
	- prefix (first 1)
	- suffix (last 3)
	- POSTag
	- wordflag
*/

//go:generate stringer -type=contextType
type contextType byte

const featuresPerContext = 8
const contexts = 5
const (
	// previous previous (prev2)
	prev2Word contextType = iota
	prev2Lemma
	prev2Cluster
	prev2Shape
	prev2Prefix1
	prev2Suffix3
	prev2POSTag
	prev2Flags

	// previous
	prevWord
	prevLemma
	prevCluster
	prevShape
	prevPrefix1
	prevSuffix3
	prevPOSTag
	prevFlags

	// ith token
	ithWord
	ithLemma
	ithCluster
	ithShape
	ithPrefix1
	ithSuffix3
	ithPOSTag
	ithFlags

	// next token
	nextWord
	nextLemma
	nextCluster
	nextShape
	nextPrefix1
	nextSuffix3
	nextPOSTag
	nextFlags

	// next next token
	next2Word
	next2Lemma
	next2Cluster
	next2Shape
	next2Prefix1
	next2Suffix3
	next2POSTag
	next2Flags

	MAXCONTEXTTYPE
)

type contextMap [MAXCONTEXTTYPE]string

func getContext(prev2, prev, ith, next, next2 *lingo.Annotation) (retVal contextMap) {
	var listOfFeats = [contexts][featuresPerContext]string{
		extractContext(prev2),
		extractContext(prev),
		extractContext(ith),
		extractContext(next),
		extractContext(next2),
	}

	for i, l := range listOfFeats {
		for j, s := range l {
			retVal[i*featuresPerContext+j] = s
		}
	}

	return retVal
}

// type featureContext struct {
// 	word    string
// 	lemma   string
// 	cluster lingo.Cluster
// 	shape   string
// 	prefix  string
// 	suffix  string
// 	POSTag  lingo.POSTag
// 	flag    lingo.WordFlag
// }

// extractContext extracts the feature contexts from a given annotation
func extractContext(a *lingo.Annotation) (retVal [featuresPerContext]string) {
	if a == nil {
		return retVal
	}

	word := a.Lowered

	// we normalize all the unicode btes first
	asRunes := []rune(a.Value)
	loweredRunes := []rune(word)

	retVal[0] = word
	retVal[1] = a.Lemma
	retVal[2] = strconv.Itoa(int(a.Cluster))
	retVal[3] = string(a.Shape)

	// prefix and suffix
	// we want the characters, not the bytes
	// for the prefix, we'll use the un-normalized version because having that extra fidelity would be useful
	if len(asRunes) > 0 {
		retVal[4] = string(asRunes[0])
	} else {
		retVal[4] = ""
	}
	if len(loweredRunes) >= 3 {
		retVal[5] = string(loweredRunes[len(loweredRunes)-3 : len(loweredRunes)])
	} else {
		retVal[5] = ""
	}
	retVal[6] = a.POSTag.String()
	retVal[7] = a.WordFlag.String()

	return retVal
}
