package dep

import (
	"testing"

	"github.com/chewxy/lingo/corpus"

	G "gorgonia.org/gorgonia"
)

func TestTrainerInitializations(t *testing.T) {
	var d *Trainer
	c := corpus.New()

	d = NewTrainer(WithCorpus(c))
	if d.corpus != c {
		t.Errorf("Expected Corpus to be set to %p. Got %p instead", c, d.corpus)
	}

	d = NewTrainer(WithConfig(DefaultNNConfig))
	if d.corpus != KnownWords {
		t.Error("Expected corpus to be set to the default KnownWords corpus")
	}
	if d.nn == nil {
		t.Fatal("Expected a neural network")
	}
	if d.nn.dict != KnownWords {
		t.Error("Expected neuralnetwork's dict to be set")
	}

	// d2 = d.Clone()
	// if d2.nn != d.nn {
	// 	t.Error("Expected a neural network!")
	// }

	// // init empty
	// d = New()
	// if err := d.Init(); err != nil {
	// 	t.Errorf("%+v", err)
	// }

	// // init with a corpus
	// d = New(WithCorpus(c))
	// if err := d.Init(); err != nil {
	// 	t.Errorf("%+v", err)
	// }
}

func TestTrainer_train(t *testing.T) {
	sts := allSentences()
	epochs := 10

	var err error

	trainer := NewTrainer(WithGeneratedCorpus(sts...), WithTrainingSet(sts))
	if err = trainer.Train(epochs); err == nil {
		t.Error("Expected an error when training an uninitialized Trainer")
	}

	// with init
	t.Logf("Pass On Costs Directly")
	conf := DefaultNNConfig
	conf.BatchSize = 90
	trainer = NewTrainer(WithGeneratedCorpus(sts...), WithConfig(conf), WithTrainingSet(sts))
	if err := trainer.Init(); err != nil {
		t.Errorf("%+v", err)
	}
	trainer.PassDirect = true

	var costs []float64
	cost := trainer.Cost()

	go func() {
		for c := range cost {
			costs = append(costs, c)
			t.Logf("Cost %v", c)
		}
	}()

	if err = trainer.Train(epochs); err != nil {
		t.Errorf("Err: %v", err)
	}

	if len(costs) == 0 {
		t.Errorf("Zero costs...")
		goto avgcosts
	}

	t.Logf("Costs %d", len(costs))
	if len(costs) < (epochs*2)-5 { // we'll allow some tolerance
		t.Errorf("Expected some costs")
	}
	if costs[0] < costs[len(costs)-1] {
		t.Errorf("Costs should be reducing")
	}

avgcosts:
	// with init, avg costs
	t.Logf("Average Costs")
	costs = costs[:0] // reset
	conf = DefaultNNConfig
	conf.Dtype = G.Float32

	trainer = NewTrainer(WithGeneratedCorpus(sts...), WithConfig(conf), WithTrainingSet(sts))
	if err := trainer.Init(); err != nil {
		t.Errorf("%+v", err)
	}
	trainer.PassDirect = false

	cost = trainer.Cost()

	go func() {
		for c := range cost {
			costs = append(costs, c)
			t.Logf("Cost %v", c)
		}
	}()
	if err = trainer.Train(epochs); err != nil {
		t.Errorf("%v", err)
	}

	if len(costs) == 0 {
		t.Fatal("Zero costs")
	}

	t.Logf("Costs %d", len(costs))
	if len(costs) == 0 {
		t.Errorf("Expected some costs")
	}

	if costs[0] < costs[len(costs)-1] {
		t.Errorf("Costs should be reducing")
	}
}

func TestTestTrainer_crossValidateTrain(t *testing.T) {
	sts := allSentences()
	cv := cvSentences()
	epochs := 10

	var trainer *Trainer
	var err error

	// uninit
	t.Logf("Uninitiated")
	trainer = NewTrainer(WithGeneratedCorpus(sts...))
	if err = trainer.Train(epochs); err == nil {
		t.Errorf("Expected an error when training with an uninitialized Trainer")
	}

	// with init
	t.Logf("Pass On Costs Directly")
	conf := DefaultNNConfig
	conf.BatchSize = 90
	trainer = NewTrainer(WithGeneratedCorpus(sts...), WithConfig(conf), WithTrainingSet(sts), WithCrossValidationSet(cv))
	trainer.PassDirect = true
	if err := trainer.Init(); err != nil {
		t.Errorf("%+v", err)
	}

	var costs []float64
	cost := trainer.Cost()
	perf := trainer.Perf()

	go func() {
		for p := range perf {
			t.Logf("Perf \n%v", p)
		}
	}()

	go func() {
		for c := range cost {
			costs = append(costs, c)
			t.Logf("Cost %v", c)
		}
	}()
	if err = trainer.Train(epochs); err != nil {
		t.Error(err)
	}

	if len(costs) == 0 {
		t.Errorf("Zero costs")
		goto avgCosts
	}

	t.Logf("Costs %d", len(costs))
	if len(costs) < (epochs*2)-5 { // we'll allow some tolerance
		t.Errorf("Expected some costs")
	}
	if costs[0] < costs[len(costs)-1] {
		t.Errorf("Costs should be reducing")
	}

avgCosts:
	// with init, avg costs, and using float32
	t.Logf("Average Costs")
	costs = costs[:0] // reset
	conf = DefaultNNConfig
	conf.Dtype = G.Float32
	trainer = NewTrainer(WithGeneratedCorpus(sts...), WithConfig(conf), WithTrainingSet(sts), WithCrossValidationSet(cv))
	if err := trainer.Init(); err != nil {
		t.Errorf("%+v", err)
	}
	trainer.PassDirect = false

	cost = trainer.Cost()
	perf = trainer.Perf()

	go func() {
		for p := range perf {
			t.Logf("Perf \n%v", p)
		}
	}()

	go func() {
		for c := range cost {
			costs = append(costs, c)
			t.Logf("Cost %v", c)
		}
	}()
	trainer.Train(epochs)

	if len(costs) == 0 {
		t.Fatal("Zero costs")
	}

	t.Logf("Costs %d", len(costs))
	if len(costs) == 0 {
		t.Errorf("Expected some costs")
	}

	if costs[0] < costs[len(costs)-1] {
		t.Errorf("Costs should be reducing")
	}
}
