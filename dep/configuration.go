package dep

import (
	"fmt"

	"github.com/chewxy/lingo"
)

// describes the current state of the parser

type head int

const (
	DOES_NOT_EXIST head = iota - 1
)

// configuration is the meat of the shift-reduce parsing. It holds the state for the shift reduction
type configuration struct {
	*lingo.Dependency
	stack  []head
	buffer []head

	bp int // buffer pointer - starts at 0, increments
}

func newConfiguration(sentence lingo.AnnotatedSentence, fromGold bool) *configuration {
	if fromGold {
		sentence = sentence.Clone()
	}

	dep := lingo.NewDependency(lingo.FromAnnotatedSentence(sentence), lingo.AllocTree())
	dep.SetID()
	sentence = sentence[1:] // because the POSTagger automatically adds a ROOTTAG at the end of it

	var buffer []head
	for i := 1; i <= len(sentence); i++ {
		buffer = append(buffer, head(i))
	}

	var stack []head
	stack = append(stack, head(0)) // add root

	return &configuration{
		Dependency: dep,
		stack:      stack,
		buffer:     buffer,
	}
}

func (c *configuration) String() string {
	return fmt.Sprintf("Stack: %v Buffer(%d): %v", c.stack, c.bp, c.buffer[c.bp:])
}

func (c *configuration) GoString() string {
	return fmt.Sprintf("Stack: %v Buffer(%d): %v\nHeads: %v\nRels: %v\n", c.stack, c.bp, c.buffer[c.bp:], c.Heads(), c.Labels())
}

func (c *configuration) bufferSize() int {
	return len(c.buffer) - c.bp
}

func (c *configuration) stackSize() int {
	return len(c.stack)
}

func (c *configuration) head(i int) head {
	heads := c.Heads() // TODO: maybe some sanity checks?
	return head(heads[i])
}

// gets the sentence index of the ith word on the stack. If there isn't anything on the stack, it returns DOES_NOT_EXIST
func (c *configuration) stackValue(i int) head {
	size := c.stackSize()
	if i >= size || i < 0 {
		return DOES_NOT_EXIST
	}
	return c.stack[size-1-i]
}

func (c *configuration) bufferValue(i int) head {
	size := c.bufferSize()
	if i >= size {
		return DOES_NOT_EXIST
	}
	return c.buffer[i+c.bp]
}

/*  stack machinations */

// pop pops the stack. It isn't really used any more. removeStack(), removeTopStack() and removeSecondTopStack() has superseded its function
func (c *configuration) pop() head {
	retVal := c.stack[len(c.stack)-1]
	c.stack = c.stack[0 : len(c.stack)-1]
	return retVal
}

// removes a value from the stack.
func (c *configuration) removeStack(i int) {
	c.stack = c.stack[:i+copy(c.stack[i:], c.stack[i+1:])]
}

// removeSecondTopStack removes the 2nd-to-last element
func (c *configuration) removeSecondTopStack() bool {
	stackSize := c.stackSize()
	if stackSize < 2 {
		return false
	}
	i := stackSize - 2
	c.removeStack(i)
	return true
}

func (c *configuration) removeTopStack() bool {
	stackSize := c.stackSize()
	if stackSize < 1 {
		return false
	}
	i := stackSize - 1
	c.removeStack(i)
	return true
}

/* Dependency related stuff */

func (c *configuration) label(i head) lingo.DependencyType {
	if i < 0 {
		return lingo.NoDepType
	}

	if i == 0 {
		return lingo.NoDepType
	}

	return c.Label(int(i))
	// i--

	// labels := c.Labels()
	// return labels[i]
}

func (c *configuration) annotation(i head) *lingo.Annotation {
	if i < 0 {
		return lingo.NullAnnotation()
	}

	if i == 0 {
		return lingo.RootAnnotation()
	}
	// i--

	return c.Annotation(int(i))

	// return c.Sentence()[i]
}

// gets the jth left child of the ith word of a sentence
func (c *configuration) lc(k, cnt head) head {
	if k < 0 || int(k) > c.N() {
		return DOES_NOT_EXIST
	}

	cc := 0
	for i := 1; i < int(k); i++ {
		if c.Head(i) == int(k) {
			cc++
			if int(cnt) == cc {
				return head(i)
			}
		}
	}
	return DOES_NOT_EXIST
}

func (c *configuration) rc(k, cnt head) head {
	if k < 0 || int(k) > c.N() {
		return DOES_NOT_EXIST
	}

	cc := 0
	for i := c.N(); i > int(k); i-- {
		if c.Head(i) == int(k) {
			cc++
			if cc == int(cnt) {
				return head(i)
			}
		}
	}
	return DOES_NOT_EXIST
}

func (c *configuration) hasOtherChildren(i int, goldParse *lingo.Dependency) bool {
	for j := 1; j <= goldParse.N(); j++ {
		if goldParse.Head(j) == i && c.Head(j) != i {
			return true
		}
	}
	return false
}

func (c *configuration) isTerminal() bool {
	return c.stackSize() == 1 && c.bufferSize() == 0
}

// Actual Transitioning stuff
func (c *configuration) shift() bool {
	i := c.bufferValue(0)
	if i == DOES_NOT_EXIST {
		return false
	}

	c.bp++ // move the buffer pointer up

	c.stack = append(c.stack, i) // push to it.... gotta work the pop
	return true
}
