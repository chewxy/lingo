package pos

import (
	"bufio"
	"encoding/gob"
	"io"
	"os"

	"github.com/chewxy/lingo"
)

// Model is the model that the POS Tagger runs on.
type Model struct {
	*perceptron
	cachedTags map[string]lingo.POSTag
}

// Save saves the model
func (m *Model) Save(filename string) error {
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

	if err := encoder.Encode(m.perceptron); err != nil {
		return err
	}

	if err := encoder.Encode(m.cachedTags); err != nil {
		return err
	}

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

	m := &Model{
		perceptron: newPerceptron(),
	}
	if err := decoder.Decode(m.perceptron); err != nil {
		return nil, err
	}

	if err := decoder.Decode(&m.cachedTags); err != nil {
		return nil, err
	}

	return m, nil

}

func (p *Tagger) Load(filename string) error {
	m, err := Load(filename)
	if err != nil {
		return err
	}
	p.Model = m
	return nil
}
