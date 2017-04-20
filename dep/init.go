package dep

import "github.com/chewxy/lingo/corpus"

func init() {
	c := corpus.New()
	c.Add("") // add null words

	KnownWords = c
}
