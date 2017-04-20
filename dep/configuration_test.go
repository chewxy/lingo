package dep

import (
	"testing"

	"github.com/chewxy/lingo"
	"github.com/stretchr/testify/assert"
)

func TestStackAppendRemove(t *testing.T) {
	sentence := mediumSentence()[0]
	as := sentence.AnnotatedSentence(dummyFix{})

	c := newConfiguration(as, true)
	t.Logf("C: %v", c)
	t.Logf("C: %#v", c)

	assert := assert.New(t)

	c.stack = append(c.stack, 200)
	assert.Equal([]head{0, 200}, c.stack, "stack is not equal after appending")

	c.removeTopStack()
	assert.Equal([]head{0}, c.stack, "stack is not equal after removeTopStack")

	c.stack = append(c.stack, 200)
	c.removeSecondTopStack()
	assert.Equal([]head{200}, c.stack, "stack is not equal after removeSecondTopStack()")

	correctHeads := []int{-1} // the -1 is the root
	correctHeads = append(correctHeads, sentence.Heads...)
	correctLabels := []lingo.DependencyType{lingo.Root}
	correctLabels = append(correctLabels, sentence.Labels...)

	dep := sentence.Dependency(dummyFix{})
	assert.Equal(correctHeads, dep.Heads(), "Heads are not equal")
	assert.Equal(correctLabels, dep.Labels(), "Labels are not equal %v \n %v", correctLabels, dep.Labels())
}

func TestConfiguration_StackValue(t *testing.T) {
	c := new(configuration)
	c.stack = []head{0, 1, 2, 5, 6}

	zero := c.stackValue(0)
	one := c.stackValue(1)
	four := c.stackValue(4)
	five := c.stackValue(5)
	negone := c.stackValue(-1)

	assert := assert.New(t)
	assert.Equal(head(6), zero, "Zeroth value not the same")
	assert.Equal(head(5), one, "First value not the same")
	assert.Equal(head(0), four, "Fourth value not the same")
	assert.Equal(DOES_NOT_EXIST, five, "Fifth value not the same")
	assert.Equal(DOES_NOT_EXIST, negone, "NegOne value not the same")

}
