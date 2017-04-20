// +build stanfordtags

package pos

import (
	"bytes"
	"encoding/gob"
	"testing"

	"github.com/chewxy/lingo"
	"github.com/stretchr/testify/assert"
)

func TestFeatureSerialization(t *testing.T) {
	var f, f2 feature
	f = singleFeature{ithWord_, "hello"}
	f2 = tupleFeature{ithWord_, "hello", "world"}

	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	decoder := gob.NewDecoder(&buf)

	if err := encoder.Encode(&f); err != nil {
		t.Fatal(err)
	}

	if err := encoder.Encode(&f2); err != nil {
		t.Fatal(err)
	}

	var decodedF, decodedF2 feature
	if err := decoder.Decode(&decodedF); err != nil {
		t.Fatal(err)
	}

	if err := decoder.Decode(&decodedF2); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, f, decodedF, "feature not deserialized properly")
	assert.Equal(t, f2, decodedF2, "feature not deserialized properly")
}

func TestPerceptron_Serialize(t *testing.T) {
	p := newPerceptron()

	// set up a dummy weight
	f := singleFeature{ithWord_, "hello"}
	w := new([lingo.MAXTAG]float64)
	w[lingo.NN] = 0.5
	w[lingo.VB] = 0.1
	p.weights[f] = w

	fc := fctuple{f, lingo.VB}
	p.totals[fc] = 0.1337
	p.steps[fc] = 0.65535

	p.instancesSeen = 1022

	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	decoder := gob.NewDecoder(&buf)

	// encode
	if err := encoder.Encode(p); err != nil {
		t.Fatal(err)
	}

	// decode
	p2 := newPerceptron()
	if err := decoder.Decode(p2); err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)

	assert.Equal(p.weights, p2.weights, "The weights have not been deserialized properly")
	assert.Equal(p.totals, p2.totals, "Totals have not been deserialized properly")
	assert.Equal(p.steps, p2.steps, "Steps have not been deserialized properly")
	assert.Equal(p.instancesSeen, p2.instancesSeen, "InstancesSeen not deserialized properly")
}
