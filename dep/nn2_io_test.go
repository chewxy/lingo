package dep

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"testing"

	"github.com/chewxy/lingo"
	"github.com/chewxy/lingo/corpus"
	G "gorgonia.org/gorgonia"
)

func TestNNIO(t *testing.T) {
	sts := allSentences()
	nn := new(neuralnetwork2)
	nn.NNConfig = DefaultNNConfig
	nn.dict = corpus.GenerateCorpus(sts)
	nn.transitions = transitions

	if err := nn.init(); err != nil {
		t.Fatalf("%+v", err)
	}

	s := `Config
------
Batch Size               : 10000
Dropout Rate             : 0.500000
AdaGrad Eps (ε)          : 0.000001
AdaGrad Learn Rate (η)   : 0.010000
Regularization Parameter : 0.000002
Hidden Layer Size        : 200
Embedding Size           : 50
Number Precomputed       : 30000

Evaluate Per 100 Iterations
Clear Gradients Per 0 Iterations
Dtype: float64

Info
------
Embeddings_Word       : (74, 50)
Embeddings_POStag     : (%d, 50)
Embeddings_Dependency : (%d, 50)
Selects_Words         : 18
Selects_POSTag        : 18
Selects_Dependency    : 12
Weights1_Word         : (200, 900)
Weights1_POSTag       : (200, 900)
Weights1_Dependency   : (200, 600)
Biases                : (200)
Weights2              : (%d, 200)
`

	correctDesc := fmt.Sprintf(s, lingo.MAXTAG, lingo.MAXDEPTYPE, MAXTRANSITION)
	if nn.String() != correctDesc {
		t.Errorf("Oops. Got %q. Want %q", nn.String(), correctDesc)
	}
	// nn.Dtype = tensor.Float32

	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(nn); err != nil {
		t.Fatalf("%+v", err)
	}

	decoder := gob.NewDecoder(&buf)
	nn2 := new(neuralnetwork2)
	nn2.dict = corpus.GenerateCorpus(sts)
	nn2.transitions = transitions
	if err := decoder.Decode(nn2); err != nil {
		t.Fatal(err)
	}

	if nn.String() != correctDesc {
		t.Fatalf("Oops. Got %q. Want %q", nn.String(), correctDesc)
	}

	if !G.ValueEq(nn.e_w.Value(), nn2.e_w.Value()) {
		t.Errorf("Expected e_w to be the same. Expected %1.1s. Got %1.1s", nn.e_w.Value(), nn2.e_w.Value())
	}

	if !G.ValueEq(nn.e_t.Value(), nn2.e_t.Value()) {
		t.Errorf("Expected e_t to be the same. Expected %1.1s. Got %1.1s", nn.e_t.Value(), nn2.e_t.Value())
	}

	if !G.ValueEq(nn.e_l.Value(), nn2.e_l.Value()) {
		t.Errorf("Expected e_l to be the same. Expected %1.1s. Got %1.1s", nn.e_l.Value(), nn2.e_l.Value())
	}

	if !G.ValueEq(nn.w1_w.Value(), nn2.w1_w.Value()) {
		t.Errorf("Expected w1_w to be the same. Expected %1.1s. Got %1.1s", nn.w1_w.Value(), nn2.w1_w.Value())
	}

	if !G.ValueEq(nn.w1_t.Value(), nn2.w1_t.Value()) {
		t.Errorf("Expected w1_t to be the same. Expected %1.1s. Got %1.1s", nn.w1_t.Value(), nn2.w1_t.Value())
	}

	if !G.ValueEq(nn.w1_l.Value(), nn2.w1_l.Value()) {
		t.Errorf("Expected w1_l to be the same. Expected %1.1s. Got %1.1s", nn.w1_l.Value(), nn2.w1_l.Value())
	}

	if !G.ValueEq(nn.b.Value(), nn2.b.Value()) {
		t.Errorf("Expected b to be the same. Expected %1.1s. Got %1.1s", nn.b.Value(), nn2.b.Value())
	}

	if !G.ValueEq(nn.w2.Value(), nn2.w2.Value()) {
		t.Errorf("Expected w2 to be the same. Expected %1.1s. Got %1.1s", nn.w2.Value(), nn2.w2.Value())
	}

	t.Logf("Visual Inspection: \n%+1.8s\n%+1.8s", nn.e_w.Value(), nn2.e_w.Value())

	// special case
	buf.Reset()
	encoder = gob.NewEncoder(&buf)
	if err := encoder.Encode(nn); err != nil {
		t.Fatalf("%+v", err)
	}
	decoder = gob.NewDecoder(&buf)
	nn3 := new(neuralnetwork2)
	if err := decoder.Decode(nn3); err == nil {
		t.Error("Expected a nocorpus error")
	}
}
