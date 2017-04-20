package dep

import (
	"fmt"

	"github.com/chewxy/lingo"
)

// transition is a tuple of Move and label
type transition struct {
	Move
	lingo.DependencyType
}

var transitions []transition
var MAXTRANSITION int

func buildTransitions(labels []lingo.DependencyType) []transition {
	ts := make([]transition, 0)
	// for _, l := range labels {
	// 	if l == lingo.NoDepType {
	// 		continue
	// 	}
	// 	t := transition{Left, l}
	// 	ts = append(ts, t)
	// }

	// for _, l := range labels {
	// 	if l == lingo.NoDepType {
	// 		continue
	// 	}

	// 	t := transition{Right, l}
	// 	ts = append(ts, t)
	// }

	// ts = append(ts, transition{Shift, lingo.NoDepType})

	for _, m := range ALLMOVES {
		for _, l := range labels {
			if (m == Shift && l != lingo.NoDepType) || (m != Shift && l == lingo.NoDepType) {
				continue
			}
			t := transition{m, l}
			ts = append(ts, t)
		}
	}
	return ts
}

func (t transition) String() string {
	return fmt.Sprintf("(%s, %s)", t.Move, t.DependencyType)
}

func lookupTransition(t transition, table []transition) int {
	for i, v := range table {
		if v == t {
			return i
		}
	}
	panic(fmt.Sprintf("Transition %v not found", t))
}

// this builds the default transitions
func init() {
	lbls := make([]lingo.DependencyType, lingo.MAXDEPTYPE)

	for i := 0; i < int(lingo.MAXDEPTYPE); i++ {
		lbls[i] = lingo.DependencyType(i)
	}

	transitions = buildTransitions(lbls)
	MAXTRANSITION = len(transitions)
}
