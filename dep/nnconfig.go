package dep

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/chewxy/gorgonia/tensor"
	"github.com/pkg/errors"
)

// NNConfig configures the neural network
type NNConfig struct {
	BatchSize                  int     // 10000
	Dropout                    float64 // 0.5
	AdaEps                     float64 // 1e-6
	AdaAlpha                   float64 //0.02
	Reg                        float64 // 1e-8
	HiddenSize                 int     // 200
	EmbeddingSize              int     // 50
	NumPrecomputed             int     //100000
	EvalPerIteration           int     // 100
	ClearGradientsPerIteration int     // 0

	Dtype tensor.Dtype
}

func (c NNConfig) String() string {
	s := `Batch Size               : %d
Dropout Rate             : %f
AdaGrad Eps (ε)          : %f
AdaGrad Learn Rate (η)   : %f
Regularization Parameter : %f
Hidden Layer Size        : %d
Embedding Size           : %d
Number Precomputed       : %d

Evaluate Per %d Iterations
Clear Gradients Per %d Iterations
Dtype: %v
`
	return fmt.Sprintf(s, c.BatchSize, c.Dropout, c.AdaEps, c.AdaAlpha, c.Reg, c.HiddenSize, c.EmbeddingSize, c.NumPrecomputed, c.EvalPerIteration, c.ClearGradientsPerIteration, c.Dtype)
}

// DefaultNNConfig is the default config that is passed in, for initialization purposses.
var DefaultNNConfig NNConfig

func (c NNConfig) GobEncode() ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	encoder.Encode(c.BatchSize)
	encoder.Encode(c.Dropout)
	encoder.Encode(c.AdaEps)
	encoder.Encode(c.AdaAlpha)
	encoder.Encode(c.Reg)
	encoder.Encode(c.HiddenSize)
	encoder.Encode(c.EmbeddingSize)
	encoder.Encode(c.NumPrecomputed)
	encoder.Encode(c.EvalPerIteration)
	encoder.Encode(c.ClearGradientsPerIteration)

	switch c.Dtype {
	case tensor.Float64:
		encoder.Encode(byte(0))
	case tensor.Float32:
		encoder.Encode(byte(1))
	default:
		return nil, errors.Errorf("Unsupported Dtype to be GobEncoded")
	}
	return buf.Bytes(), nil
}

func (c *NNConfig) GobDecode(p []byte) error {
	b := bytes.NewBuffer(p)
	decoder := gob.NewDecoder(b)

	decoder.Decode(&c.BatchSize)
	decoder.Decode(&c.Dropout)
	decoder.Decode(&c.AdaEps)
	decoder.Decode(&c.AdaAlpha)
	decoder.Decode(&c.Reg)
	decoder.Decode(&c.HiddenSize)
	decoder.Decode(&c.EmbeddingSize)
	decoder.Decode(&c.NumPrecomputed)
	decoder.Decode(&c.EvalPerIteration)
	decoder.Decode(&c.ClearGradientsPerIteration)

	var bite byte
	decoder.Decode(&bite)
	switch bite {
	case 0:
		c.Dtype = tensor.Float64
	case 1:
		c.Dtype = tensor.Float32
	default:
		return errors.Errorf("Unsupported Dtype to be GobDecoded: %v", bite)
	}
	return nil
}

func init() {
	DefaultNNConfig = NNConfig{
		BatchSize: 10000,
		Dropout:   0.5,

		AdaEps:   1e-6,
		AdaAlpha: 0.01,

		Reg: 1.5e-6,

		HiddenSize:     200,
		EmbeddingSize:  50,
		NumPrecomputed: 30000,

		EvalPerIteration:           100,
		ClearGradientsPerIteration: 0,

		Dtype: tensor.Float64,
		// Dtype: gorgonia.Float32,
	}
}
