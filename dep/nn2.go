package dep

import (
	"github.com/chewxy/lingo"
	"github.com/chewxy/lingo/corpus"
	"github.com/pkg/errors"
	G "gorgonia.org/gorgonia"
	"gorgonia.org/tensor"
)

// may is a simple monad for handling errors
type may struct {
	error
	n *G.Node
}

func (m *may) doUnary(fn func(*G.Node) (*G.Node, error)) {
	if m.error != nil {
		return
	}
	m.n, m.error = fn(m.n)
}

func (m *may) doBinary(fn func(a, b *G.Node) (*G.Node, error), other *G.Node) {
	if m.error != nil {
		return
	}
	m.n, m.error = fn(m.n, other)
}

func (m *may) doSwapBinary(fn func(a, b *G.Node) (*G.Node, error), other *G.Node) {
	if m.error != nil {
		return
	}
	m.n, m.error = fn(other, m.n)
}

type neuralnetwork2 struct {
	NNConfig

	g   *G.ExprGraph
	sub *G.ExprGraph

	// model

	// embedding matrices for word, POSTags and labels respectively
	e_w *G.Node // Shape: (EmbeddingSize, DictSize)
	e_t *G.Node // Shape: (EmbeddingSize, lingo.MAXTAG)
	e_l *G.Node // Shape: (EmbeddingSize, lingo.MAXDEP)

	// w1
	w1_w *G.Node // Shape: (HiddenSize, DictSize)
	w1_t *G.Node // Shape: (HiddenSize, lingo.MAXTAG)
	w1_l *G.Node // Shape: (HiddenSize, lingo.MAXDEP)
	b    *G.Node // Shape: (HiddenSize)

	// w2
	w2 *G.Node // Shape: (MAXTRANSITION, HiddenSize)

	// selects
	x_wSelW G.Nodes // 18 - word features
	x_tSelT G.Nodes // 18 - POSTag features
	x_lSelL G.Nodes // 12 - Dependency feature

	// inputs (feature vectors built up from the selects)
	x_w *G.Node
	x_t *G.Node
	x_l *G.Node

	// outputs
	scores  *G.Node // argmax this to get the greedy decoded transition
	logProb *G.Node
	cost    *G.Node
	costVal G.Value

	vm     G.VM
	model  G.Nodes
	solver G.Solver

	dict        *corpus.Corpus
	transitions []transition

	costChan chan G.Value

	// wordfeats *G.Node
	// tagfeats  *G.Node
	// depfeats  *G.Node
	// sumfeats  *G.Node
	// act       *G.Node
}

func (nn *neuralnetwork2) initialized() bool {
	return nn.g != nil && nn.sub != nil &&
		nn.e_w != nil && nn.e_t != nil && nn.e_l != nil &&
		nn.w1_w != nil && nn.w1_t != nil && nn.w1_l != nil && nn.b != nil &&
		nn.w2 != nil && len(nn.x_wSelW) > 0 && len(nn.x_tSelT) > 0 && len(nn.x_lSelL) > 0 &&
		nn.x_w != nil && nn.x_t != nil && nn.x_l != nil &&
		nn.scores != nil &&
		nn.dict != nil && nn.vm != nil && nn.solver != nil
}

func (nn *neuralnetwork2) init() error {
	if nn.dict == nil {
		return errors.Errorf("No Corpus Provided to the Neural Network. Will be unable to decode")
	}

	g := G.NewGraph()
	nn.g = g

	word := nn.dict.Size()
	tags := int(lingo.MAXTAG)
	deps := int(lingo.MAXDEPTYPE)
	// trns := len(nn.transitions)

	wordFeats := POS_OFFSET - 0
	tagFeats := DEP_OFFSET - POS_OFFSET
	depFeats := int(MAXFEATURE) - DEP_OFFSET

	// In any case a very very very small dict was passed in
	// we set the minimum to wordFeatss
	if word < wordFeats {
		word = wordFeats
	}

	logf(`Word: %d
tags: %d
deps: %d
wordFeats: %d
tagFeats: %d
depFeats: %d
`, word, tags, deps, wordFeats, tagFeats, depFeats)

	// define inputs
	nn.x_w = G.NewVector(g, nn.Dtype, G.WithShape(wordFeats*nn.EmbeddingSize), G.WithName("word input"), G.WithInit(G.Zeroes()))
	nn.x_t = G.NewVector(g, nn.Dtype, G.WithShape(tagFeats*nn.EmbeddingSize), G.WithName("POSTag input"), G.WithInit(G.Zeroes()))
	nn.x_l = G.NewVector(g, nn.Dtype, G.WithShape(depFeats*nn.EmbeddingSize), G.WithName("word input"), G.WithInit(G.Zeroes()))

	nn.x_wSelW = make(G.Nodes, wordFeats)
	nn.x_tSelT = make(G.Nodes, tagFeats)
	nn.x_lSelL = make(G.Nodes, depFeats)

	// define models
	nn.e_w = G.NewMatrix(g, nn.Dtype, G.WithShape(word, nn.EmbeddingSize), G.WithName("e_w"), G.WithInit(G.GlorotU(1)))
	nn.e_t = G.NewMatrix(g, nn.Dtype, G.WithShape(tags, nn.EmbeddingSize), G.WithName("e_t"), G.WithInit(G.GlorotU(1)))
	nn.e_l = G.NewMatrix(g, nn.Dtype, G.WithShape(deps, nn.EmbeddingSize), G.WithName("e_l"), G.WithInit(G.GlorotU(1)))

	nn.w1_w = G.NewMatrix(g, nn.Dtype, G.WithShape(nn.HiddenSize, nn.EmbeddingSize*wordFeats), G.WithName("w1_w"), G.WithInit(G.GlorotU(1)))
	nn.w1_t = G.NewMatrix(g, nn.Dtype, G.WithShape(nn.HiddenSize, nn.EmbeddingSize*tagFeats), G.WithName("w1_t"), G.WithInit(G.GlorotU(1)))
	nn.w1_l = G.NewMatrix(g, nn.Dtype, G.WithShape(nn.HiddenSize, nn.EmbeddingSize*depFeats), G.WithName("w1_l"), G.WithInit(G.GlorotU(1)))
	nn.b = G.NewVector(g, nn.Dtype, G.WithShape(nn.HiddenSize), G.WithName("b"), G.WithInit(G.Zeroes()))

	nn.w2 = G.NewMatrix(g, nn.Dtype, G.WithShape(MAXTRANSITION, nn.HiddenSize), G.WithName("w2"), G.WithInit(G.GlorotU(1)))

	nn.model = G.Nodes{nn.e_w, nn.e_t, nn.e_l, nn.w1_w, nn.w1_t, nn.w1_l, nn.b, nn.w2}

	// define selects
	// words first
	logf("nn.e_w: %+1.1s", nn.e_w.Value())
	var err error
	for i := 0; i < wordFeats; i++ {
		if nn.x_wSelW[i], err = G.Slice(nn.e_w, G.S(i)); err != nil { // dummy slices... they'll be replaced at runtime
			return err
		}

	}

	// tag features
	for i := 0; i < tagFeats; i++ {
		if nn.x_tSelT[i], err = G.Slice(nn.e_t, G.S(i)); err != nil { // dummy slices... they'll be replaced at runtime
			return err
		}
	}

	// dependency features
	for i := 0; i < depFeats; i++ {
		if nn.x_lSelL[i], err = G.Slice(nn.e_l, G.S(i)); err != nil {
			return err
		}
	}

	// forwards
	if err = nn.fwd(); err != nil {
		return err
	}

	// backprop
	if _, err = G.Grad(nn.cost, nn.model...); err != nil {
		return err
	}

	nn.sub = g.SubgraphRoots(nn.scores)

	// prog, locmap, err := G.Compile(nn.g)
	// if err != nil {
	// 	return err
	// }
	// log.Printf("Prog: %v", prog)

	// ioutil.WriteFile("graph.dot", []byte(g.ToDot()), 0644)

	// logger := log.New(os.Stderr, "", 0)
	// nn.vm = G.NewTapeMachine(prog, locmap, G.BindDualValues(nn.model...), G.UseCudaFor(), G.WithLogger(logger), G.WithWatchlist())
	// nn.vm = G.NewTapeMachine(prog, locmap, G.BindDualValues(nn.model...), G.UseCudaFor())
	nn.vm = G.NewTapeMachine(nn.g, G.BindDualValues(nn.model...), G.UseCudaFor())
	G.BindDualValues(nn.scores)(nn.vm) // makes sure that scores is a *dualValue
	nn.solver = G.NewAdaGradSolver(G.WithLearnRate(nn.AdaAlpha), G.WithEps(nn.AdaEps), G.WithL2Reg(nn.Reg), G.WithBatchSize(float64(nn.BatchSize)))
	// nn.solver = G.NewVanillaSolver(G.WithLearnRate(nn.AdaAlpha), G.WithL2Reg(nn.Reg))
	return nil
}

func (nn *neuralnetwork2) fwd() error {
	var err error

	// build up x vectors
	if nn.x_w, err = G.Concat(0, nn.x_wSelW...); err != nil {
		return err
	}

	if nn.x_t, err = G.Concat(0, nn.x_tSelT...); err != nil {
		return err
	}

	if nn.x_l, err = G.Concat(0, nn.x_lSelL...); err != nil {
		return err
	}

	logf("w1_w %v, x_w %v", nn.w1_w.Shape(), nn.x_w.Shape())
	m_w := &may{nil, nn.w1_w}
	m_w.doBinary(G.Mul, nn.x_w)
	if m_w.error != nil {
		return m_w.error
	}

	logf("w1_t %v, x_t %v", nn.w1_t.Shape(), nn.x_t.Shape())
	m_t := &may{nil, nn.w1_t}
	m_t.doBinary(G.Mul, nn.x_t)
	if m_t.error != nil {
		return m_t.error
	}

	logf("w1_l %v, x_l %v", nn.w1_l.Shape(), nn.x_l.Shape())
	m_l := &may{nil, nn.w1_l}
	m_l.doBinary(G.Mul, nn.x_l)
	if m_l.error != nil {
		return m_l.error
	}

	// add and activate layer 1
	logf("w : %v", m_w.n.Shape())
	m_w1 := &may{nil, m_w.n}
	m_w1.doBinary(G.Add, m_t.n)
	m_w1.doBinary(G.Add, m_l.n)
	m_w1.doBinary(G.Add, nn.b)
	m_w1.doUnary(G.Cube)
	if m_w1.error != nil {
		return m_w1.error
	}

	if nn.Dropout > 0 {
		logf("Doing dropout")
		m_w1.n, m_w1.error = G.Dropout(m_w1.n, nn.Dropout)
		if m_w1.error != nil {
			return m_w1.error
		}
	}

	// go to softmax layer
	logf("w2: %v, w1act: %v", nn.w2.Shape(), m_w1.n.Shape())
	m_sm := &may{nil, nn.w2}
	m_sm.doBinary(G.Mul, m_w1.n)
	nn.scores = m_sm.n
	m_sm.doUnary(G.SoftMax)
	if m_sm.error != nil {
		return m_sm
	}

	nn.logProb = m_sm.n
	// G.WithName("Logprob")(nn.logProb)
	// log.Printf("LOGPROB %v %p %v", nn.logProb, nn.logProb, nn.logProb)
	if nn.cost, err = G.Slice(nn.logProb, G.S(0)); err != nil { // slice is a dummy tensor.Slice. It'll be replaced at runtime
		return err
	}

	G.Read(nn.cost, &nn.costVal)
	return nil
}

func (nn *neuralnetwork2) costProgress() <-chan G.Value {
	if nn.costChan == nil {
		nn.costChan = make(chan G.Value)
	}
	return nn.costChan
}

// train does one epoch of training. The examples are batched.
func (nn *neuralnetwork2) train(examples []example) error {
	size := len(examples)
	batches := size / nn.BatchSize

	var start, end int
	if nn.BatchSize > size {
		batches = 1
		end = size
		G.WithBatchSize(float64(size))(nn.solver) // set it such that the solver doesn't get confused
	} else {
		end = nn.BatchSize
	}

	for batch := 0; batch < batches; batch++ {
		for _, ex := range examples[start:end] {
			nn.feats2vec(ex.features)
			tid := lookupTransition(ex.transition, nn.transitions)

			if err := G.UnsafeLet(nn.cost, G.S(tid)); err != nil {
				return err
			}

			if err := nn.vm.RunAll(); err != nil {
				return err
			}

			nn.vm.Reset()
		}
		if err := nn.solver.Step(G.NodesToValueGrads(nn.model)); err != nil {
			err = errors.Wrapf(err, "Stepping on the model failed %v", batch)
			return err
		}

		if nn.costChan != nil {
			nn.costChan <- nn.costVal
		}

		start = end
		if start >= size {
			break
		}
		end += nn.BatchSize
		if end >= size {
			end = size
		}
	}

	return nil
}

// pred predicts the index of the transitions
func (nn *neuralnetwork2) pred(ind []int) (int, error) {
	nn.feats2vec(ind)

	// f, _ := os.OpenFile("LOOOOOG", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	// logger := log.New(f, "", 0)
	// logger := log.New(os.Stderr, "", 0)

	// m := G.NewLispMachine(nn.sub, G.ExecuteFwdOnly(), G.WithLogger(logger), G.WithWatchlist(), G.LogBothDir(), G.WithValueFmt("%+3.3v"))
	m := G.NewLispMachine(nn.sub, G.ExecuteFwdOnly())
	if err := m.RunAll(); err != nil {
		return 0, err
	}
	// logger.Println("========================\n")

	val := nn.scores.Value().(tensor.Tensor)
	t, err := tensor.Argmax(val, tensor.AllAxes)
	if err != nil {
		return 0, err
	}

	return t.ScalarValue().(int), nil
}

// utility function

func (nn *neuralnetwork2) feats2vec(indicators []int) error {
	// fix word features
	for i, ind := range indicators[:POS_OFFSET] {
		if err := G.UnsafeLet(nn.x_wSelW[i], G.S(ind-wordFeatsStartAt)); err != nil {
			return err
		}
	}

	// fix tag features
	for i, ind := range indicators[POS_OFFSET:DEP_OFFSET] {
		if err := G.UnsafeLet(nn.x_tSelT[i], G.S(ind)); err != nil {
			return err
		}
	}

	for i, ind := range indicators[DEP_OFFSET:] {
		if err := G.UnsafeLet(nn.x_lSelL[i], G.S(ind-labelFeatsStartAt)); err != nil {
			return err
		}
	}

	return nil
}
