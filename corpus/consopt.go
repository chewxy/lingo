package corpus

import (
	"sync/atomic"

	"github.com/xtgo/set"
)

// ConsOpt is a construction option for manual creation of a Corpus
type ConsOpt func(c *Corpus) error

// WithWords creates a corpus from a
func WithWords(a []string) ConsOpt {
	f := func(c *Corpus) error {
		s := set.Strings(a)
		c.words = s
		c.frequencies = make([]int, len(s))

		ids := make(map[string]int)
		maxID := len(s)
		totalFreq := len(s)
		var maxWL int
		for i, w := range s {
			if len(w) > maxWL {
				maxWL = len(w)
			}
			ids[w] = i
		}
		for _, w := range s {
			c.frequencies[ids[w]]++
		}

		c.ids = ids
		atomic.AddInt64(&c.maxid, int64(maxID))
		c.totalFreq = totalFreq
		c.maxWordLength = maxWL
		return nil
	}
	return f
}

// WithOrderedWords creates a Corpus with the given word order
func WithOrderedWords(a []string) ConsOpt {
	f := func(c *Corpus) error {
		s := a
		c.words = s
		c.frequencies = make([]int, len(s))
		for i := range c.frequencies {
			c.frequencies[i] = 1
		}

		ids := make(map[string]int)
		maxID := len(s)
		totalFreq := len(s)
		var maxWL int
		for i, w := range a {
			if len(w) > maxWL {
				maxWL = len(w)
			}
			ids[w] = i
		}

		c.ids = ids
		atomic.AddInt64(&c.maxid, int64(maxID))
		c.totalFreq = totalFreq
		c.maxWordLength = maxWL
		return nil
	}
	return f
}

// WithSize preallocates all the things in Corpus
func WithSize(size int) ConsOpt {
	return func(c *Corpus) error {
		c.words = make([]string, 0, size)
		c.frequencies = make([]int, 0, size)
		return nil
	}
}
