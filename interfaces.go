package lingo

import (
	"encoding/gob"

	"github.com/chewxy/gorgonia/tensor"
)

// Lemmatizer is anything that can lemmatize
type Lemmatizer interface {
	Lemmatize(string, POSTag) ([]string, error)
}

// Stemmer is anything that can stem
type Stemmer interface {
	Stem(string) (string, error)
}

// Sentencer is anything that returns an AnnotatedSentence
type Sentencer interface {
	Sentence() AnnotatedSentence
}

// Corpus is the interface for the corpus.
type Corpus interface {
	// ID returns the ID of a word and whether or not it was found in the corpus
	Id(word string) (id int, ok bool)

	// Word returns the word given the ID, and whether or not it was found in the corpus
	Word(id int) (word string, ok bool)

	// Add adds a word to the corpus and returns its ID. If a word was previously in the corpus, it merely updates the frequency count and returns the ID
	Add(word string) int

	// Size returns the size of the corpus.
	Size() int

	// WordFreq returns the frequency of the word. If the word wasn't in the corpus, it returns 0.
	WordFreq(word string) int

	// IDFreq returns the frequency of a word given an ID. If the word isn't in the corpus it returns 0.
	IDFreq(id int) int

	// TotalFreq returns the total number of words ever seen by the corpus. This number includes the count of repeat words.
	TotalFreq() int

	// MaxWordLength returns the length of the longest known word in the corpus
	MaxWordLength() int

	// WordProb returns the probability of a word appearing in the corpus
	WordProb(word string) (float64, bool)

	// IO stuff
	gob.GobEncoder
	gob.GobDecoder
}

// WordEmbeddings is any type that is both a corpus and can return word vectors
type WordEmbeddings interface {
	Corpus

	// WordVector returns a vector of embeddings given the word
	WordVector(word string) (vec tensor.Tensor, err error)

	// Vector returns a vector of embeddings given the word ID
	Vector(id int) (vec tensor.Tensor, err error)

	// Embedding returns the matrix
	Embedding() tensor.Tensor
}
