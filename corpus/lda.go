package corpus

import (
	"github.com/chewxy/gorgonia"
	"gorgonia.org/tensor"
)

// LDAModel ... TODO
//https://en.wikipedia.org/wiki/Latent_Dirichlet_allocation
type LDAModel struct {
	// params
	Alpha tensor.Tensor   // is a Row
	Eta   tensor.Tensor   // is a Col
	Kappa gorgonia.Scalar // Decay
	Tau0  gorgonia.Scalar // offset

	// parameters needed for working
	Topics      int
	ChunkSize   int
	Terms       int
	UpdateEvery int
	EvalEvery   int

	// consts
	Iterations     int
	GammaThreshold float64

	MinimumProb float64

	// track current progress
	Updates int

	// type
	Dtype tensor.Dtype
}

func (l *LDAModel) init() {
	eta := tensor.New(tensor.Of(l.Dtype), tensor.WithShape(l.Topics))
	alpha := tensor.New(tensor.Of(l.Dtype), tensor.WithShape(l.Topics))

	switch l.Dtype {
	case tensor.Float64:
		v := 1.0 / float64(l.Topics)
		eta.Memset(v)
		alpha.Memset(v)
	case tensor.Float32:
		v := float32(1) / float32(l.Topics)
		eta.Memset(v)
		alpha.Memset(v)
	}

	l.Alpha = alpha
	l.Eta = eta
}
