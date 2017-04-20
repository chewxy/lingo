package pos

import (
	"math"

	"github.com/chewxy/lingo"
)

func maxScore(scores *[lingo.MAXTAG]float64) lingo.POSTag {
	var maxClass lingo.POSTag
	maxVal := -math.MaxFloat64
	for c, v := range scores {
		if v > maxVal {
			maxClass = lingo.POSTag(c)
			maxVal = v
		}
	}

	return maxClass
}
