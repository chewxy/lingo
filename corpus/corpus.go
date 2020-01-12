package corpus

import (
	"errors"
	"sync/atomic"
	"unicode/utf8"
)

// Corpus is a data structure holding the relevant metadata and information for a corpus of text.
// It serves as vocabulary with ID for lookup. This is very useful as neural networks rely on the IDs rather than the text themselves
type Corpus struct {
	words       []string
	frequencies []int

	ids map[string]int

	// atomic read and write plz
	maxid         int64
	totalFreq     int
	maxWordLength int
}

// New creates a new *Corpus
func New() *Corpus {
	c := &Corpus{
		words:       make([]string, 0),
		frequencies: make([]int, 0),
		ids:         make(map[string]int),
	}

	// add some default words
	c.Add("") // aka NULL - when there are no words
	c.Add("-UNKNOWN-")
	c.Add("-ROOT-")
	c.maxWordLength = 0 // specials don't have lengths

	return c
}

// Construct creates a Corpus given the construction options. This allows for more flexibility
func Construct(opts ...ConsOpt) (*Corpus, error) {
	c := new(Corpus)

	// checks
	if c.words == nil {
		c.words = make([]string, 0)
	}
	if c.frequencies == nil {
		c.frequencies = make([]int, 0)
	}
	if c.ids == nil {
		c.ids = make(map[string]int)
	}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

// ID returns the ID of a word and whether or not it was found in the corpus
func (c *Corpus) Id(word string) (int, bool) {
	id, ok := c.ids[word]
	return id, ok
}

// Word returns the word given the ID, and whether or not it was found in the corpus
func (c *Corpus) Word(id int) (string, bool) {
	size := atomic.LoadInt64(&c.maxid)
	maxid := int(size)

	if id >= maxid {
		return "", false
	}
	return c.words[id], true
}

// Add adds a word to the corpus and returns its ID. If a word was previously in the corpus, it merely updates the frequency count and returns the ID
func (c *Corpus) Add(word string) int {
	if id, ok := c.ids[word]; ok {
		c.frequencies[id]++
		c.totalFreq++
		return id
	}

	id := atomic.AddInt64(&c.maxid, 1)
	c.ids[word] = int(id - 1)
	c.words = append(c.words, word)
	c.frequencies = append(c.frequencies, 1)
	c.totalFreq++

	runeCount := utf8.RuneCountInString(word)
	if runeCount > c.maxWordLength {
		c.maxWordLength = runeCount
	}

	return int(id - 1)
}

// Size returns the size of the corpus.
func (c *Corpus) Size() int {
	size := atomic.LoadInt64(&c.maxid)
	return int(size)
}

// WordFreq returns the frequency of the word. If the word wasn't in the corpus, it returns 0.
func (c *Corpus) WordFreq(word string) int {
	id, ok := c.ids[word]
	if !ok {
		return 0
	}

	return c.frequencies[id]
}

// IDFreq returns the frequency of a word given an ID. If the word isn't in the corpus it returns 0.
func (c *Corpus) IDFreq(id int) int {
	size := atomic.LoadInt64(&c.maxid)
	maxid := int(size)

	if id >= maxid {
		return 0
	}
	return c.frequencies[id]
}

// TotalFreq returns the total number of words ever seen by the corpus. This number includes the count of repeat words.
func (c *Corpus) TotalFreq() int {
	return c.totalFreq
}

// MaxWordLength returns the length of the longest known word in the corpus.
func (c *Corpus) MaxWordLength() int {
	return c.maxWordLength
}

// WordProb returns the probability of a word appearing in the corpus.
func (c *Corpus) WordProb(word string) (float64, bool) {
	id, ok := c.Id(word)
	if !ok {
		return 0, false
	}

	count := c.frequencies[id]
	return float64(count) / float64(c.totalFreq), true

}

// Merge combines two corpuses. The receiver is the one that is mutated.
func (c *Corpus) Merge(other *Corpus) {
	for i, word := range other.words {
		freq := other.frequencies[i]
		if id, ok := c.ids[word]; ok {
			c.frequencies[id] += freq
			c.totalFreq += freq
		} else {
			id := c.Add(word)
			c.frequencies[id] += freq - 1
			c.totalFreq += freq - 1
		}
	}
}

// Replace replaces the content of a word. The old reference remains.
//
// e.g: c.Replace("foo", "bar")
// c.Id("foo") will still return a ID. The ID will be the same as c.Id("bar")
func (c *Corpus) Replace(a, with string) error {
	old, ok := c.ids[a]
	if !ok {
		return errors.Errorf("Cannot replace %q with %q. %q is not found", a, with, a)
	}
	if _, ok := c.ids[with]; ok {
		return errors.Errorf("Cannot replace %q with %q. %q exists in the corpus", a, with, with)
	}
	c.words[old] = with
	return nil

}

func (c *Corpus) ReplaceWord(id int, with string) error {
	if id >= len(c.words) {
		return errors.Errorf("ID %d out of bounds", id)
	}
	c.words[id] = with
	return nil
}
