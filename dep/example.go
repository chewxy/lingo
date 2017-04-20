package dep

import (
	"math/rand"

	"github.com/chewxy/lingo"
	"github.com/chewxy/lingo/corpus"
	"github.com/chewxy/lingo/treebank"
)

// example is a training example.
type example struct {
	transition

	features []int // features are used in the embeddings
	labels   []int // labels are used in scoring the transitions
}

func makeExamples(sentenceTags []treebank.SentenceTag, conf NNConfig, dict *corpus.Corpus, ts []transition, f lingo.AnnotationFixer) []example {
	var examples []example

	var tarpit, nonprojective, good int
	for i, sentenceTag := range sentenceTags {
		exs, err := makeOneExample(i, sentenceTag, dict, ts, f)
		if err != nil {
			switch err.(type) {
			case TarpitError:
				tarpit++
			case NonProjectiveError:
				nonprojective++
			}
		} else {
			examples = append(examples, exs...)
			good++
		}
	}

	logf("Number of SentenceTags Generated Into Examples: %d/%d | Number of Examples: %d | Number of nonprojective examples: %d | Number of tarpit examples: %d", good, len(sentenceTags), len(examples), nonprojective, tarpit)
	return examples
}

// makeOneExample is an example of a poorly named function. It makes an example from a SentenceTag
func makeOneExample(i int, sentenceTag treebank.SentenceTag, dict *corpus.Corpus, ts []transition, f lingo.AnnotationFixer) ([]example, error) {
	var examples []example

	s := sentenceTag.AnnotatedSentence(f)
	dep := s.Dependency()
	if dep.IsProjective() {
		c := newConfiguration(s, true)

		count := 0
		for !c.isTerminal() && count < 1000 {
			if count == 999 {
				return examples, TarpitError{c}
			}

			oracle := c.oracle(dep)
			features := getFeatures(c, dict)

			labels := make([]int, MAXTRANSITION)
			for i, t := range ts {
				if t == oracle {
					labels[i] = 1
				} else if c.canApply(t) {
					labels[i] = 0
				} else {
					labels[i] = -1
				}
			}

			ex := example{transition{oracle.Move, oracle.DependencyType}, features, labels}
			examples = append(examples, ex)

			c.apply(oracle)
			count++
		}
	} else {
		return nil, NonProjectiveError{dep}
	}

	return examples, nil
}

func shuffleExamples(a []example) {
	for i := range a {
		j := rand.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}
}
