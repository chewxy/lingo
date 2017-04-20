package dep

import (
	"fmt"
	"os"
	"sync"

	"github.com/chewxy/lingo"
	"github.com/chewxy/lingo/corpus"
	"github.com/chewxy/lingo/treebank"
	"github.com/pkg/errors"
)

// TrainerConsOpt is a construction option for trainer
type TrainerConsOpt func(t *Trainer)

// WithTrainingModel loads a trainer with a model
func WithTrainingModel(m *Model) TrainerConsOpt {
	f := func(t *Trainer) {
		t.Model = m
	}
	return f
}

// WithTrainingSet creates a trainer with a training set
func WithTrainingSet(st []treebank.SentenceTag) TrainerConsOpt {
	f := func(t *Trainer) {
		t.trainingSet = st
	}
	return f
}

// WithCrossValidationSet creates a trainer with a cross validation set
func WithCrossValidationSet(st []treebank.SentenceTag) TrainerConsOpt {
	f := func(t *Trainer) {
		t.crossValSet = st
	}
	return f
}

// WithConfig sets up a *Trainer with a NNConfig
func WithConfig(conf NNConfig) TrainerConsOpt {
	f := func(t *Trainer) {
		t.nn.NNConfig = conf
		t.nn.dict = t.corpus
		t.nn.transitions = t.ts
		t.EvalPerIter = conf.EvalPerIteration
	}
	return f
}

// WithLemmatizer sets the lemmatizer option on the Trainer
func WithLemmatizer(l lingo.Lemmatizer) TrainerConsOpt {
	f := func(t *Trainer) {
		// cannot pass in itself!
		if T, ok := l.(*Trainer); ok && T == t {
			panic("Recursive definition of lemmatizer (trying to set the t.lemmatizer = T) !")
		}

		t.l = l
	}
	return f
}

// WithStemmer sets up the stemmer option on the DependencyParser
func WithStemmer(s lingo.Stemmer) TrainerConsOpt {
	f := func(t *Trainer) {
		// cannot pass in itself
		if T, ok := s.(*Trainer); ok && T == t {
			panic("Recursive setting of stemmer! (Trying to set t.stemmer = T)")
		}
		t.s = s
	}
	return f
}

// WithCluster sets the brown cluster options for the DependencyParser
func WithCluster(c map[string]lingo.Cluster) TrainerConsOpt {
	f := func(t *Trainer) {
		t.c = c
	}
	return f
}

// WithCorpus creates a Trainer with a corpus
func WithCorpus(c *corpus.Corpus) TrainerConsOpt {
	f := func(t *Trainer) {
		t.corpus = c
		t.nn.dict = c
	}
	return f
}

// WithGeneratedCorpus creates a Trainer's corpus from a list of SentenceTags. The corpus will be generated from the SentenceTags
func WithGeneratedCorpus(sts ...treebank.SentenceTag) TrainerConsOpt {
	f := func(t *Trainer) {
		dict := corpus.GenerateCorpus(sts)
		if t.corpus == nil {
			t.corpus = dict
		} else {
			t.corpus.Merge(dict)
		}

		t.nn.dict = t.corpus
	}
	return f
}

// Trainer trains a model
type Trainer struct {
	trainingSet []treebank.SentenceTag
	crossValSet []treebank.SentenceTag

	once sync.Once
	*Model

	// Training configuration
	EvalPerIter int    // for cross validation - evaluate results every n epochs
	PassDirect  bool   // Pass on the costs directly to the cost channel? If false, an average will be used
	SaveBest    string // SaveBest is the filename that will be saved. If it's empty then the best-while-training will not be saved

	// fixer
	l lingo.Lemmatizer
	s lingo.Stemmer
	c map[string]lingo.Cluster

	err  chan error
	cost chan float64
	perf chan Performance
}

// NewTrainer creates a new Trainer.
func NewTrainer(opts ...TrainerConsOpt) *Trainer {
	t := new(Trainer)
	// set up the default model
	t.Model = new(Model)
	t.corpus = KnownWords
	t.ts = transitions

	// set up the neural network
	t.nn = new(neuralnetwork2)
	t.nn.NNConfig = DefaultNNConfig
	t.nn.transitions = transitions
	t.nn.dict = KnownWords

	for _, opt := range opts {
		opt(t)
	}
	return t
}

// Lemmatize implemnets lingo.Lemmatizer
func (t *Trainer) Lemmatize(a string, pt lingo.POSTag) ([]string, error) {
	if t.l == nil {
		return nil, componentUnavailable("Lemmatizer")
	}
	return t.l.Lemmatize(a, pt)
}

// Stem implements lingo.Stemmer
func (t *Trainer) Stem(a string) (string, error) {
	if t.s == nil {
		return "", componentUnavailable("Stemmer")
	}
	return t.s.Stem(a)
}

// Clusters implements lingo.Fixer
func (t *Trainer) Clusters() (map[string]lingo.Cluster, error) {
	if t.c == nil {
		return nil, componentUnavailable("Clusters")
	}
	return t.c, nil
}

/* Getters */

// Cost returns a channel of costs for monitoring the training. If the PassDirect field in the trainer is set to true
// then the costs are directly returned. Otherwise the costs are averaged over the epoch.
func (t *Trainer) Cost() <-chan float64 {
	if t.cost == nil {
		t.cost = make(chan float64)
	}
	return t.cost
}

// Perf returns a channel of Performance for monitoring the training.
func (t *Trainer) Perf() <-chan Performance {
	if t.perf == nil {
		t.perf = make(chan Performance)
	}
	return t.perf
}

/* Methods */

// Init initializes the DependencyParser with a corpus and a neural network config
func (t *Trainer) Init() (err error) {
	f := func() {
		err = t.nn.init()
	}
	t.once.Do(f)
	return
}

// Train trains a model.
//
// If a cross validation set is provided, it will automatically train with the cross validation set
func (t *Trainer) Train(epochs int) error {
	if err := t.pretrainCheck(); err != nil {
		return err
	}
	if len(t.crossValSet) > 0 {
		return t.crossValidateTrain(epochs)
	}
	return t.train(epochs)
}

// TrainWithoutCrossValidation trains a model without cross validation.
func (t *Trainer) TrainWithoutCrossValidation(epochs int) error {
	return t.train(epochs)
}

// train simply trains the model without having a cross validation.
func (t *Trainer) train(epochs int) error {

	var epochChan chan struct{}
	if t.cost != nil {
		defer func() {
			close(t.cost)
			t.cost = nil
		}()

		epochChan = t.handleCosts()
		if epochChan != nil {
			defer close(epochChan)
		}
	}

	examples := makeExamples(t.trainingSet, t.nn.NNConfig, t.nn.dict, t.ts, t)

	for e := 0; e < epochs; e++ {
		if err := t.nn.train(examples); err != nil {
			return err
		}

		if epochChan != nil {
			epochChan <- struct{}{}
		}

		shuffleExamples(examples)
	}
	return nil
}

// crossValidateTrain trains the model but also does cross validation to ensure overfitting don't happen.
func (t *Trainer) crossValidateTrain(epochs int) error {
	if t.perf != nil {
		defer func() {
			close(t.perf)
			t.perf = nil
		}()
	}

	var epochChan chan struct{}
	if t.cost != nil {
		defer func() {
			close(t.cost)
			t.cost = nil
		}()

		epochChan = t.handleCosts()
		if epochChan != nil {
			defer close(epochChan)
		}
	}
	examples := makeExamples(t.trainingSet, t.nn.NNConfig, t.nn.dict, t.ts, t)

	var best Performance
	for e := 0; e < epochs; e++ {
		if err := t.nn.train(examples); err != nil {
			return err
		}

		if t.EvalPerIter > 0 && e%t.EvalPerIter == 0 || e == epochs-1 {
			perf := t.crossValidate(t.crossValSet)

			// if there is a channel to report back the performance, send it down
			if t.perf != nil {
				perf.Iter = e
				t.perf <- perf
			}

			if perf.UAS > best.UAS {
				best = perf

				if t.SaveBest != "" {
					f, err := os.Create(t.SaveBest)
					if err != nil {
						err = errors.Wrapf(err, "Unable to open SaveBest file %q", t.SaveBest)
						return err
					}

					t.Model.SaveWriter(f)
				}
			}
		}

		if epochChan != nil {
			epochChan <- struct{}{}
		}

		shuffleExamples(examples)
	}
	return nil
}

// pretrainCheck checks if everything is sane
func (t *Trainer) pretrainCheck() error {
	// check
	if t.nn == nil || !t.nn.initialized() {
		return errors.Errorf("DependencyParser not init()'d. Perhaps you forgot to call .Init() somewhere?")
	}

	if len(t.trainingSet) == 0 {
		return errors.Errorf("Cannot train with no training data set")
	}

	return nil
}

// handleCosts handles the costs from the neural network in two ways:
//		1. pass: directly passes on the costs (which may come from multiple batches in an epoch)
//		2. mean: calculates the mean of the costs and passes it on into d.cost
//
// If d.cost is nil, it simply returns. This method should be called after a check that d.cost is not nil
func (t *Trainer) handleCosts() (epochChan chan struct{}) {
	nncost := t.nn.costProgress()

	if t.PassDirect {
		go func() {
			for cost := range nncost {
				switch c := cost.Data().(type) {
				case float32:
					t.cost <- float64(c)
				case float64:
					t.cost <- c
				default:
					// this should NEVER happen
					panic(fmt.Sprintf("Unhandled cost type %T", c))
				}
			}
		}()
	} else {
		epochChan = make(chan struct{})

		// it collects the costs until the epoch chan signals that an epoch is done. Then the cost is averaged and sent down the d.cost channel
		go func(epochChan chan struct{}) {
			var collected []float64
			for {
				select {
				case cost := <-nncost:
					switch c := cost.Data().(type) {
					case float32:
						collected = append(collected, float64(c))
					case float64:
						collected = append(collected, c)
					default:
						// this should NEVER happen
						panic(fmt.Sprintf("Unhandled cost type %T", c))
					}
				case <-epochChan:
					var avg float64
					for _, cost := range collected {
						avg += cost
					}

					if len(collected) > 0 {
						avg /= float64(len(collected))
					}

					t.cost <- avg
					collected = collected[:0]
				}
			}
		}(epochChan)
	}
	return
}
