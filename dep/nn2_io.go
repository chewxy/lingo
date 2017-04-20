package dep

import (
	"bytes"
	"encoding/gob"
	"fmt"

	G "github.com/chewxy/gorgonia"
	T "github.com/chewxy/gorgonia/tensor"
	"github.com/pkg/errors"
)

var empty struct{}

func (nn *neuralnetwork2) String() string {
	s := `Config
------
%v
Info
------
Embeddings_Word       : %v
Embeddings_POStag     : %v
Embeddings_Dependency : %v
Selects_Words         : %d
Selects_POSTag        : %d
Selects_Dependency    : %d
Weights1_Word         : %v
Weights1_POSTag       : %v
Weights1_Dependency   : %v
Biases                : %v
Weights2              : %v
`

	return fmt.Sprintf(s, nn.NNConfig,
		nn.e_w.Shape(), nn.e_t.Shape(), nn.e_l.Shape(),
		len(nn.x_wSelW), len(nn.x_tSelT), len(nn.x_lSelL),
		nn.w1_w.Shape(), nn.w1_t.Shape(), nn.w1_l.Shape(),
		nn.b.Shape(), nn.w2.Shape())
}

func (nn *neuralnetwork2) GobEncode() ([]byte, error) {
	if !nn.initialized() {
		return nil, errors.Errorf("Neural network not initialized. Cannot gob")
	}

	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)

	if err := encoder.Encode(nn.NNConfig); err != nil {
		return nil, err
	}

	if err := encoder.Encode(nn.e_w.Value()); err != nil {
		return nil, err
	}

	if err := encoder.Encode(nn.e_t.Value()); err != nil {
		return nil, err
	}

	if err := encoder.Encode(nn.e_l.Value()); err != nil {
		return nil, err
	}

	if err := encoder.Encode(nn.w1_w.Value()); err != nil {
		return nil, err
	}

	if err := encoder.Encode(nn.w1_t.Value()); err != nil {
		return nil, err
	}

	if err := encoder.Encode(nn.w1_l.Value()); err != nil {
		return nil, err
	}

	if err := encoder.Encode(nn.b.Value()); err != nil {
		return nil, err
	}

	if err := encoder.Encode(nn.w2.Value()); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (nn *neuralnetwork2) GobDecode(buf []byte) error {
	// prechecks
	if nn.dict == nil {
		return errors.Errorf("Neural Network has no corpus attached to it (Corpuses are serialized separately).")
	}

	b := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(b)

	if err := decoder.Decode(&nn.NNConfig); err != nil {
		return err
	}

	if err := nn.init(); err != nil {
		return err
	}

	e_w := T.New(T.Of(nn.Dtype), T.WithShape(nn.e_w.Shape()...))
	if err := decoder.Decode(e_w); err != nil {
		return err
	}
	G.Let(nn.e_w, e_w)

	e_t := T.New(T.Of(nn.Dtype), T.WithShape(nn.e_t.Shape()...))
	if err := decoder.Decode(e_t); err != nil {
		return err
	}
	G.Let(nn.e_t, e_t)

	e_l := T.New(T.Of(nn.Dtype), T.WithShape(nn.e_l.Shape()...))
	if err := decoder.Decode(e_l); err != nil {
		return err
	}
	G.Let(nn.e_l, e_l)

	w1_w := T.New(T.Of(nn.Dtype), T.WithShape(nn.w1_w.Shape()...))
	if err := decoder.Decode(w1_w); err != nil {
		return err
	}
	G.Let(nn.w1_w, w1_w)

	w1_t := T.New(T.Of(nn.Dtype), T.WithShape(nn.w1_t.Shape()...))
	if err := decoder.Decode(w1_t); err != nil {
		return err
	}
	G.Let(nn.w1_t, w1_t)

	w1_l := T.New(T.Of(nn.Dtype), T.WithShape(nn.w1_l.Shape()...))
	if err := decoder.Decode(w1_l); err != nil {
		return err
	}
	G.Let(nn.w1_l, w1_l)

	bias := T.New(T.Of(nn.Dtype), T.WithShape(nn.b.Shape()...))
	if err := decoder.Decode(bias); err != nil {
		return err
	}
	G.Let(nn.b, bias)

	w2 := T.New(T.Of(nn.Dtype), T.WithShape(nn.w2.Shape()...))
	if err := decoder.Decode(w2); err != nil {
		return err
	}
	G.Let(nn.w2, w2)

	return nil
}
