package dep

import (
	"testing"

	"github.com/chewxy/lingo/corpus"
)

func TestMakeExamples(t *testing.T) {
	st := simpleSentence()
	dict := corpus.GenerateCorpus(st)

	exs := makeExamples(st, DefaultNNConfig, dict, transitions, dummyFix{})
	if len(exs) != 20 {
		t.Error("Expected 20 examples to be generated from simple sentence")
	}
}
