package dep

import (
	"fmt"
	"log"

	"github.com/chewxy/lingo"
	"github.com/chewxy/lingo/corpus"
	"github.com/pkg/errors"
)

var KnownWords *corpus.Corpus // package provided global

// Parser is the object that performs the dependency parsing
// It contains a neural network, which is the core of it.
//
// The same object can be used to train the NN
type Parser struct {
	Input  chan lingo.AnnotatedSentence
	Output chan *lingo.Dependency
	Error  chan error

	*Model
}

// New creates a new Parser
func New(m *Model) *Parser {
	d := &Parser{
		Output: make(chan *lingo.Dependency),
		Error:  make(chan error),

		Model: m,
	}

	return d
}

// Run is used when using the NN to parse a sentence. For training, see Train()
func (d *Parser) Run() {
	defer close(d.Output)
	for sentence := range d.Input {
		log.Printf("Sentence: %d", len(sentence))
		dep, err := d.predict(sentence)

		if err != nil {
			d.Error <- err
			return
		}
		d.Output <- dep
	}
	return
}

func (d *Parser) predict(sentence lingo.AnnotatedSentence) (*lingo.Dependency, error) {
	c := newConfiguration(sentence, false)

	var err error
	var argmax int
	var count int
	for !c.isTerminal() && count < 100 {
		logf("%v", c)
		if count == 99 {
			logf("TARPIT")
		}

		features := getFeatures(c, d.corpus)
		// features2 := getFeatureArray(c, d.dict)

		if argmax, err = d.nn.pred(features); err != nil {
			return nil, err
		}
		// log.Printf("Argmax: %v, len(d.ts): %v, len(transitions) %v", argmax, len(d.ts), len(transitions))
		t := transitions[argmax] // no this is NOT a mistake
		if !c.canApply(t) {
			t = transition{Shift, lingo.NoDepType} // reset
			// manual argmaxing
			switch scores := d.nn.scores.Value().Data().(type) {
			case []float32:
				var maxScore float32
				for i, kt := range d.ts {
					if scores[i] > maxScore && c.canApply(kt) {
						maxScore = scores[i]
						t = kt
					}
				}
			case []float64:
				var maxScore float64
				for i, kt := range d.ts {
					if scores[i] > maxScore && c.canApply(kt) {
						maxScore = scores[i]
						t = kt
					}
				}
			default:
				return nil, errors.Errorf("Unhandled score type %T", d.nn.scores.Value())
			}

		}
		c.apply(t)

		count++
	}
	fix(c.Dependency)
	return c.Dependency, err
}

func (d *Parser) String() string {
	var nns, ds string

	if d.corpus != nil {
		ds = fmt.Sprintf("\nDict Size: %d words\nMAXTAG: %d\nMAXDEPTYPE: %d\n", d.corpus.Size(), lingo.MAXTAG, lingo.MAXDEPTYPE)
	} else {
		ds = "\n"
	}

	if d.nn != nil && d.nn.initialized() {
		nns = fmt.Sprintf("\nNeural Network:\n=================\n%v\n", d.nn)
	}

	if !d.nn.initialized() {
		panic(fmt.Sprintf("%v", d.nn))
	}

	base := "\n\nDependency Parser Info:\n=======================\n"
	return base + ds + nns
}
