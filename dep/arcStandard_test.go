package dep

import (
	"testing"

	"github.com/chewxy/lingo"
	"github.com/stretchr/testify/assert"
)

func TestCanApply(t *testing.T) {
	dep := simpleSentence()[0].Dependency(dummyFix{})

	buffer := make([]head, 0)
	for i := 1; i < dep.WordCount(); i++ {
		buffer = append(buffer, head(i))
	}

	stack := []head{0}

	c := &configuration{
		Dependency: dep,
		stack:      stack,
		buffer:     buffer,
	}

	assert := assert.New(t)

	logf("Start config: \n%v", c)

	rootLeft := c.canApply(transition{Left, lingo.Root})
	rootRight := c.canApply(transition{Right, lingo.Root})
	NSubjLeft := c.canApply(transition{Left, lingo.NSubj})
	NSubjRight := c.canApply(transition{Right, lingo.NSubj})
	ShiftDep := c.canApply(transition{Shift, lingo.NoDepType})

	assert.Equal(false, rootLeft, "rootLeft should be false")
	assert.Equal(false, rootRight, "rootRight should be false")
	assert.Equal(false, NSubjLeft, "NSubjLeft should be false")
	assert.Equal(false, NSubjRight, "NSubjRight should be false")
	assert.Equal(true, ShiftDep, "ShiftDep should be true")

	logf("rootRight: %v, rootLeft: %v", rootLeft, rootRight)
	logf("NSubjRight: %v, NSubjLeft: %v", NSubjRight, NSubjLeft)
	logf("ShiftDep: %v", ShiftDep)

	c.shift()
	c.shift()
	logf("%v", c)

	rootLeft = c.canApply(transition{Left, lingo.Root})
	rootRight = c.canApply(transition{Right, lingo.Root})
	NSubjLeft = c.canApply(transition{Left, lingo.NSubj})
	NSubjRight = c.canApply(transition{Right, lingo.NSubj})
	ShiftDep = c.canApply(transition{Shift, lingo.NoDepType})

	assert.Equal(true, rootLeft, "rootLeft should be true")
	assert.Equal(true, rootRight, "rootRight should be true")
	assert.Equal(true, NSubjLeft, "NSubjLeft should be true")
	assert.Equal(true, NSubjRight, "NSubjRight should be true")
	assert.Equal(true, ShiftDep, "ShiftDep should be true")

	logf("rootRight: %v, rootLeft: %v", rootLeft, rootRight)
	logf("NSubjRight: %v, NSubjLeft: %v", NSubjRight, NSubjLeft)
	logf("ShiftDep: %v", ShiftDep)
}

func TestOracle(t *testing.T) {
	st := simpleSentence()[0]
	s := st.AnnotatedSentence(nil)
	c := newConfiguration(s, true)
	d := s.Dependency()

	for count := 0; !c.isTerminal() && count < 100; count++ {
		oracle := c.oracle(d)

		if !c.canApply(oracle) && (oracle != transition{Right, lingo.Root}) {
			t.Errorf("Cannot apply %v", oracle)
			break
		}

		c.apply(oracle)
	}

	assert.Equal(t, d.Heads(), c.Heads())
}
