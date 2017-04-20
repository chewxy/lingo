package pos

import (
	"os"
	"strings"
	"testing"

	"github.com/chewxy/lingo/treebank"
	"github.com/stretchr/testify/assert"
)

func TestSaveLoad(t *testing.T) {
	pt := New()
	sentences := treebank.ReadConllu(strings.NewReader(conllu))

	pt.Train(sentences, 5)
	pt.Save("test.dat")

	pt2 := New()
	if err := pt2.Load("test.dat"); err != nil {
		os.Remove("test.dat")
		t.Fatal(err)
	}

	assert := assert.New(t)

	assert.Equal(pt.perceptron, pt2.perceptron, "POSTaggers' perceptrons are different:%p %p", pt.perceptron, pt2.perceptron)
	assert.Equal(pt.cachedTags, pt2.cachedTags, "POSTaggers' cachedTags are different")

	// cleanup
	os.Remove("test.dat")
}
