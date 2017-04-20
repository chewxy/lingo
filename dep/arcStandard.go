package dep

import "github.com/chewxy/lingo"

// var SingleRoot bool = true // make this part of a build process

// canApply checks if a particular transition can be applied
func (c *configuration) canApply(t transition) bool {

	var h head
	if t.Move == Left || t.Move == Right {
		if t.Move == Left {
			h = c.stackValue(0)
		} else {
			h = c.stackValue(1)
		}

		if h < 0 {
			return false
		}
		if h == 0 && t.DependencyType != lingo.Root {
			return false
		}
	}

	stackSize := c.stackSize()
	bufferSize := c.bufferSize()

	if t.Move == Left {
		return stackSize > 2
	}

	if t.Move == Right {
		return stackSize > 2 || (stackSize == 2 && bufferSize == 0)

		// if not single root build
		// return stackSize >= 2
	}

	return bufferSize > 0 // strange other thing...

}

// apply applies the transition
func (c *configuration) apply(t transition) {
	logf("Applying %v", t)
	w1 := int(c.stackValue(1))
	w2 := int(c.stackValue(0))

	if t.Move == Left {
		c.AddArc(w2, w1, t.DependencyType)
		c.removeSecondTopStack()
	} else if t.Move == Right {
		c.AddArc(w1, w2, t.DependencyType)
		c.removeTopStack()
	} else {
		c.shift()
	}
}

// oracle gets the gold transition given the state
func (c *configuration) oracle(goldParse *lingo.Dependency) (t transition) {
	w1 := int(c.stackValue(1))
	w2 := int(c.stackValue(0))

	if w1 > 0 && goldParse.Head(w1) == w2 {
		t.Move = Left
		t.DependencyType = goldParse.Label(w1)
		return
	} else if w1 >= 0 && goldParse.Head(w2) == w1 && !c.hasOtherChildren(w2, goldParse) {
		t.Move = Right
		t.DependencyType = goldParse.Label(w2)

		return
	}
	return // default transition is Shift
}
