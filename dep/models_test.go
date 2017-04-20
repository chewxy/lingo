package dep

import (
	"os"
	"testing"

	G "github.com/chewxy/gorgonia"
	"github.com/stretchr/testify/assert"
)

func TestModel_SaveLoad(t *testing.T) {
	assert := assert.New(t)

	testFileName := "TestSave.dat"
	m := new(Model)

	// dumb shit
	if err := m.Save(testFileName); err == nil {
		t.Error("Expected an error")
	}

	conf := DefaultNNConfig
	conf.Dtype = G.Float32
	m = new(Model)
	m.ts = transitions
	m.corpus = KnownWords

	m.nn = new(neuralnetwork2)
	m.nn.NNConfig = conf
	m.nn.dict = m.corpus

	if err := m.nn.init(); err != nil {
		t.Error(err)
	}

	if err := m.Save(testFileName); err != nil {
		t.Fatal(err)
	}

	var m2 *Model
	var err error
	if m2, err = Load(testFileName); err != nil {
		t.Error(err)

	}

	assert.Equal(m.corpus, m2.corpus, "Both Dependency Parsers need to have the same dict")

	if !G.ValueEq(m.nn.w2.Value(), m2.nn.w2.Value()) {
		t.Errorf("Expected w2 to be equal")
	}
	if !G.ValueEq(m.nn.e_w.Value(), m2.nn.e_w.Value()) {
		t.Errorf("Expected e_w to be equal")
	}

	// cleanup
	if err := os.Remove(testFileName); err != nil {
		t.Error(err)
	}
}
