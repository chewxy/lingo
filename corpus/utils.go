package corpus

import (
	"errors"
	"math"
)

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func dot(a, b []float64) (float64, error) {
	if len(a) != len(b) {
		return 0, errors.New("Differing lengths!")
	}

	var retVal float64
	for i, v := range a {
		retVal += v * b[i]
	}
	return retVal, nil
}

func mag(a []float64) (float64, error) {
	dotProd, err := dot(a, a)
	if err != nil {
		return dotProd, err
	}
	return math.Sqrt(dotProd), nil
}
