package pos

import (
	"math"
	"math/rand"
	"testing"

	"github.com/chewxy/lingo"
)

func TestMaxScore(t *testing.T) {
	rand.Seed(1337)
	scores := new([lingo.MAXTAG]float64)

	for i := range scores {
		scores[i] = rand.Float64()
		if lingo.POSTag(i) == lingo.ROOT_TAG {
			scores[i] = math.MaxFloat64
		}
	}

	tag := maxScore(scores)
	if tag != lingo.ROOT_TAG {
		t.Errorf("Expected Score #10 to be the max. Got %d instead", tag)
	}
}
