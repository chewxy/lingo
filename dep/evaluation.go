package dep

import (
	"fmt"
	"io/ioutil"

	"github.com/chewxy/lingo"
	"github.com/chewxy/lingo/treebank"
)

// Performance is a tuple that holds performance information from a training session
type Performance struct {
	Iter int     // which training iteration is this?
	UAS  float64 // Unlabelled Attachment Score
	LAS  float64 // Labeled Attachment Score
	UEM  float64 // Unlabelled Exact Match
	Root float64 // Correct Roots Ratio
}

func (p Performance) String() string {
	s := `EPO: %d
UAS: %.5f
LAS: %.5f
UEM: %.5f
ROO: %.5f`

	return fmt.Sprintf(s, p.Iter, p.UAS, p.LAS, p.UEM, p.Root)
}

// performance evaluation related code goes here

// Evaluate compares predicted trees with the gold standard trees and returns a Performance. It panics if the number of predicted trees and the number of gold trees aren't the same
func Evaluate(predictedTrees, goldTrees []*lingo.Dependency) Performance {
	if len(predictedTrees) != len(goldTrees) {
		panic(fmt.Sprintf("%d predicted trees; %d gold trees. Unable to compare", len(predictedTrees), len(goldTrees)))
	}

	var correctLabels, correctHeads, correctTrees, correctRoot, sumArcs float64
	var check int

	for i, tr := range predictedTrees {
		gTr := goldTrees[i]

		if len(tr.AnnotatedSentence) != len(gTr.AnnotatedSentence) {
			sumArcs += float64(gTr.N())

			// log.Printf("WARNING: %q and %q do not have the same length", tr, gTr)
			continue
		}

		var nCorrectHead int
		for j, a := range tr.AnnotatedSentence[1:] {
			b := gTr.AnnotatedSentence[j+1]
			if a.HeadID() == b.HeadID() {
				correctHeads++
				nCorrectHead++
			}

			if a.DependencyType == b.DependencyType {
				correctLabels++
			}
			sumArcs++
		}
		if nCorrectHead == gTr.N() {
			correctTrees++
		}
		if tr.Root() == gTr.Root() {
			correctRoot++
		}

		// check 5 per iteration
		if check < 5 {
			logf("predictedHeads: \n%v\n%v\n", tr.Heads(), gTr.Heads())
			logf("Ns: %v | %v || Correct: %v", tr.N(), gTr.N(), nCorrectHead)
			check++
		}
	}

	uas := correctHeads / sumArcs
	las := correctLabels / sumArcs
	uem := correctTrees / float64(len(predictedTrees))
	roo := correctRoot / float64(len(predictedTrees))

	return Performance{UAS: uas, LAS: las, UEM: uem, Root: roo}
}

func (t *Trainer) crossValidate(st []treebank.SentenceTag) Performance {
	preds := t.predMany(st)
	golds := make([]*lingo.Dependency, len(st))

	for i, s := range st {
		golds[i] = s.Dependency(t)
	}
	return Evaluate(preds, golds)
}

func (t *Trainer) predMany(sentenceTags []treebank.SentenceTag) []*lingo.Dependency {
	retVal := make([]*lingo.Dependency, len(sentenceTags))
	for i, st := range sentenceTags {
		dep, err := t.pred(st.AnnotatedSentence(t))
		if err != nil {
			ioutil.WriteFile("fullGraph.dot", []byte(t.nn.g.ToDot()), 0644)
			panic(fmt.Sprintf("%+v", err))
		}
		retVal[i] = dep
	}
	return retVal
}

func (t *Trainer) pred(as lingo.AnnotatedSentence) (*lingo.Dependency, error) {
	d := new(Parser)
	d.Model = t.Model

	return d.predict(as)
}
