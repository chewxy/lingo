package dep

import (
	"github.com/chewxy/lingo"
	"github.com/chewxy/lingo/corpus"
)

// getFeatures extracts the IDs to pass into the neural network. These IDs are used in the network to construct the  input layers
func getFeatures(c *configuration, dict *corpus.Corpus) []int {
	// logf("CONFIG: %v", c)
	wordFeats := make([]int, 0)
	posFeats := make([]lingo.POSTag, 0)
	labelFeats := make([]lingo.DependencyType, 0)
	unknownID, _ := dict.Id("-UNKNOWN-")

	for j := 2; j >= 0; j-- {
		index := c.stackValue(j)
		mor := c.annotation(index)

		if wordID, ok := dict.Id(mor.Value); ok {
			wordFeats = append(wordFeats, wordID)
		} else {
			wordFeats = append(wordFeats, unknownID)
		}
		posFeats = append(posFeats, mor.POSTag)
	}

	// logf("wordFeats: %v", wordFeats)

	for j := 0; j <= 2; j++ {
		index := c.bufferValue(j)
		mor := c.annotation(index)
		// logf("Want: %v Index: %d. Morpheme: %v", j, index, mor)

		if wordID, ok := dict.Id(mor.Value); ok {
			wordFeats = append(wordFeats, wordID)
		} else {
			wordFeats = append(wordFeats, unknownID)
		}
		posFeats = append(posFeats, mor.POSTag)
	}
	// logf("wordFeats: %v", wordFeats)

	for j := 0; j <= 1; j++ {
		k := c.stackValue(j)

		index := c.lc(k, 1)
		mor := c.annotation(index)
		if wordID, ok := dict.Id(mor.Value); ok {
			wordFeats = append(wordFeats, wordID)
		} else {
			wordFeats = append(wordFeats, unknownID)
		}
		posFeats = append(posFeats, mor.POSTag)
		labelFeats = append(labelFeats, c.label(index))

		index = c.rc(k, 1)
		mor = c.annotation(index)
		if wordID, ok := dict.Id(mor.Value); ok {
			wordFeats = append(wordFeats, wordID)
		} else {
			wordFeats = append(wordFeats, unknownID)
		}
		posFeats = append(posFeats, mor.POSTag)
		labelFeats = append(labelFeats, c.label(index))

		index = c.lc(k, 2)
		mor = c.annotation(index)
		if wordID, ok := dict.Id(mor.Value); ok {
			wordFeats = append(wordFeats, wordID)
		} else {
			wordFeats = append(wordFeats, unknownID)
		}
		posFeats = append(posFeats, mor.POSTag)
		labelFeats = append(labelFeats, c.label(index))

		index = c.rc(k, 2)
		mor = c.annotation(index)
		if wordID, ok := dict.Id(mor.Value); ok {
			wordFeats = append(wordFeats, wordID)
		} else {
			wordFeats = append(wordFeats, unknownID)
		}
		posFeats = append(posFeats, mor.POSTag)
		labelFeats = append(labelFeats, c.label(index))

		leftChild := c.lc(k, 1)
		index = c.lc(leftChild, 1)
		mor = c.annotation(index)
		if wordID, ok := dict.Id(mor.Value); ok {
			wordFeats = append(wordFeats, wordID)
		} else {
			wordFeats = append(wordFeats, unknownID)
		}
		posFeats = append(posFeats, mor.POSTag)
		labelFeats = append(labelFeats, c.label(index))

		rightChild := c.rc(k, 1)
		index = c.rc(rightChild, 1)
		mor = c.annotation(index)
		if wordID, ok := dict.Id(mor.Value); ok {
			wordFeats = append(wordFeats, wordID)
		} else {
			wordFeats = append(wordFeats, unknownID)
		}
		posFeats = append(posFeats, mor.POSTag)
		labelFeats = append(labelFeats, c.label(index))
	}

	// the embedding matrix is arranged thus:
	/*
		POSTag0 0, 1, ... 50
		POSTag1
		...
		MAXTAG-1
		DepType0
		DepType1
		...
		MAXDEPTYPE-1
		WordID0
		...
		WordIDN
	*/

	features := make([]int, MAXFEATURE)

	for i, w := range wordFeats {
		features[i] = w + wordFeatsStartAt
	}
	for i, t := range posFeats {
		features[i+POS_OFFSET] = int(t)
	}
	for i, l := range labelFeats {
		features[i+DEP_OFFSET] = int(l) + labelFeatsStartAt
	}

	return features
}

const (
	POS_OFFSET   int = 18
	DEP_OFFSET       = 36
	STACK_OFFSET     = 6
	STACK_NUMBER     = 6
)
