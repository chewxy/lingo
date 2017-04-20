package dep

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"os"

	"github.com/chewxy/lingo/corpus"
	"github.com/pkg/errors"
)

// Model holds the neural network that a DependencyParser uses. To train, use a Trainer
type Model struct {
	nn     *neuralnetwork2
	corpus *corpus.Corpus
	ts     []transition
}

func (m *Model) Corpus() *corpus.Corpus { return m.corpus }

func (m *Model) String() string {
	var buf bytes.Buffer
	buf.WriteString(m.nn.String())
	buf.WriteString("Transitions: [")
	for _, t := range m.ts {
		fmt.Fprintf(&buf, "%v, ", t)
	}
	buf.WriteString("]")
	return buf.String()
}

func (m *Model) Save(filename string) error {
	if m.nn == nil {
		return errors.Errorf("Cannot save a model with no nn")
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	return m.SaveWriter(f)
}

func (m *Model) SaveWriter(f io.WriteCloser) error {
	defer f.Close()
	w := bufio.NewWriter(f)
	defer w.Flush()
	encoder := gob.NewEncoder(w)

	if err := encoder.Encode(m.corpus); err != nil {
		return err
	}

	if err := encoder.Encode(m.nn); err != nil {
		return err
	}

	// if err := encoder.Encode(m.ts); err != nil {
	// 	return err
	// }

	return nil
}

func Load(filename string) (*Model, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return LoadReader(f)
}

func LoadReader(rd io.ReadCloser) (*Model, error) {
	defer rd.Close()
	r := bufio.NewReader(rd)
	decoder := gob.NewDecoder(r)

	m := new(Model)
	if err := decoder.Decode(&m.corpus); err != nil {
		return nil, err
	}

	m.nn = new(neuralnetwork2)
	m.nn.dict = m.corpus

	if err := decoder.Decode(&m.nn); err != nil {
		return nil, err
	}

	if err := decoder.Decode(&m.ts); err != nil {
		m.ts = transitions
	}
	m.nn.transitions = m.ts

	return m, nil

}
