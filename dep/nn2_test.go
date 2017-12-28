package dep

import (
	"math/rand"
	"testing"
	"time"

	"github.com/chewxy/lingo/corpus"
	"gorgonia.org/gorgonia"
)

func TestNN2(t *testing.T) {
	rand.Seed(1337)

	// we test 50 iterations unless the short flag is passed in
	epochs := 50
	if testing.Short() {
		epochs = 10
	}

	sts := allSentences()
	nn := new(neuralnetwork2)
	nn.NNConfig = DefaultNNConfig
	nn.Dtype = gorgonia.Float32
	nn.dict = corpus.GenerateCorpus(sts)
	nn.transitions = transitions

	if err := nn.init(); err != nil {
		t.Fatalf("%+v", err)
	}

	var costs []float64
	ch := nn.costProgress()
	sigChan := make(chan struct{})

	go func(ch <-chan gorgonia.Value, sig chan struct{}) {
		for cost := range ch {
			switch c := cost.Data().(type) {
			case float32:
				costs = append(costs, float64(c))
			case float64:
				costs = append(costs, c)
			}

			t.Logf("Cost %v", cost)
		}
		sig <- struct{}{}
	}(ch, sigChan)

	exs := makeExamples(sts, nn.NNConfig, nn.dict, transitions, dummyFix{})

	start := time.Now()
	for i := 0; i < epochs; i++ {
		if err := nn.train(exs); err != nil {
			t.Errorf("%+v", err)
		}
		shuffleExamples(exs)
	}
	// simulate what *DependencyParser would do
	close(nn.costChan)
	nn.costChan = nil

	t.Logf("Training %d iterations took Taken: %v", epochs, time.Since(start))

	<-sigChan
	if len(costs) == 0 {
		t.Error("Expected some costs")
	}
	if costs[0] <= costs[len(costs)-1] {
		t.Error("Expected costs to have reduced during training")
	}

	// PREDICTION TIME!

	ss2 := simpleSentence()
	exs = makeExamples(ss2, nn.NNConfig, nn.dict, transitions, dummyFix{})
	start = time.Now()
	for i, ex := range exs {
		ind, err := nn.pred(ex.features)
		if err != nil {
			t.Errorf("Example %d failed: %v", i, err)
			continue
		}

		t.Logf("Example %d. Want: %v. Got %v. Same: %t", i, ex.transition, transitions[ind], ex.transition == transitions[ind])
	}
	t.Logf("Pred Time Taken: %v", time.Since(start))
}
