package pos

import "github.com/chewxy/lingo"

type perceptron struct {
	// weights map[feature]*[lingo.MAXTAG]float64 // it's a pointer to a static array because map values are immutable, and cannot be edited

	weightsSF map[singleFeature]*[lingo.MAXTAG]float64
	weightsTF map[tupleFeature]*[lingo.MAXTAG]float64

	totals map[fctuple]float64
	steps  map[fctuple]float64

	instancesSeen float64
}

// feature-class tuple is a tuple that contains a feature and a class. This makes calculation of the averaging easier
type fctuple struct {
	feature
	lingo.POSTag
}

func newPerceptron() *perceptron {
	return &perceptron{
		// weights: make(map[feature]*[lingo.MAXTAG]float64),

		weightsSF: make(map[singleFeature]*[lingo.MAXTAG]float64),
		weightsTF: make(map[tupleFeature]*[lingo.MAXTAG]float64),

		totals: make(map[fctuple]float64),
		steps:  make(map[fctuple]float64),
	}
}

func (p *perceptron) updateWeightsSF(f singleFeature, tag lingo.POSTag, weight, value float64) {
	tuple := fctuple{f, tag}
	p.totals[tuple] += (p.instancesSeen - p.steps[tuple]) * weight
	p.steps[tuple] = p.instancesSeen

	if _, ok := p.weightsSF[f]; !ok {
		p.weightsSF[f] = new([lingo.MAXTAG]float64)
	}
	p.weightsSF[f][tag] = weight + value
}

func (p *perceptron) updateWeightsTF(f tupleFeature, tag lingo.POSTag, weight, value float64) {
	tuple := fctuple{f, tag}
	p.totals[tuple] += (p.instancesSeen - p.steps[tuple]) * weight
	p.steps[tuple] = p.instancesSeen

	if _, ok := p.weightsTF[f]; !ok {
		p.weightsTF[f] = new([lingo.MAXTAG]float64)
	}
	p.weightsTF[f][tag] = weight + value
}

func (p *perceptron) update(guess, truth lingo.POSTag, sf sfFeatures, tf tfFeatures) {
	p.instancesSeen++
	if truth == guess {
		return
	}

	for _, f := range sf {
		var truthValue float64
		var guessValue float64

		if weights, ok := p.weightsSF[f]; ok {
			truthValue = weights[truth]
			guessValue = weights[guess]
		}

		p.updateWeightsSF(f, truth, truthValue, 1)
		p.updateWeightsSF(f, guess, guessValue, -1)
	}

	for _, f := range tf {
		var truthValue float64
		var guessValue float64

		if weights, ok := p.weightsTF[f]; ok {
			truthValue = weights[truth]
			guessValue = weights[guess]
		}

		p.updateWeightsTF(f, truth, truthValue, 1)
		p.updateWeightsTF(f, guess, guessValue, -1)
	}
}

func (p *perceptron) predict(sf sfFeatures, tf tfFeatures) lingo.POSTag {
	var scores [lingo.MAXTAG]float64
	for _, f := range sf {
		if weights, ok := p.weightsSF[f]; ok {
			for label, weight := range weights {
				scores[label] += weight
			}
		}
	}

	for _, f := range tf {
		if weights, ok := p.weightsTF[f]; ok {
			for label, weight := range weights {
				scores[label] += weight
			}
		}
	}

	return maxScore(&scores)
}

func (p *perceptron) average() {
	for f, weights := range p.weightsSF {
		for c, weight := range weights {
			tuple := fctuple{f, lingo.POSTag(c)}
			total := p.totals[tuple]

			total += (p.instancesSeen - p.steps[tuple]) * weight
			avg := total / p.instancesSeen

			weights[c] = avg
		}
	}

	for f, weights := range p.weightsTF {
		for c, weight := range weights {
			tuple := fctuple{f, lingo.POSTag(c)}
			total := p.totals[tuple]

			total += (p.instancesSeen - p.steps[tuple]) * weight
			avg := total / p.instancesSeen

			weights[c] = avg
		}
	}
}
