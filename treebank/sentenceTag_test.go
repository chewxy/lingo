package treebank

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSentenceTag(t *testing.T) {
	assert := assert.New(t)
	readr := strings.NewReader(sampleConllu)
	st := ReadConllu(readr)[0]

	correctHeads := []int{2, 5, 4, 5, 0, 7, 5, 9, 5, 11, 9, 14, 14, 11, 18, 18, 18, 14, 5}
	assert.Equal(correctHeads, st.Heads)

	dep := st.Dependency(nil)
	assert.Equal(correctHeads, dep.Heads()[1:])
}
