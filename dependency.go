package lingo

import (
	"bytes"
	"fmt"
)

// Dependency represents the dependency parse of a sentence. While AnnotatedSentence does
// already do a job of representing the dependency parse of a sentence, *Dependency actually contains
// meta information about the dependency parse (specifically, lefts, rights) that makes parsing a dependency a lot faster
//
// The fields are mostly left unexported for a good reason - a dependency parse SHOULD be static after it's been built
type Dependency struct {
	AnnotatedSentence

	wordCount int

	lefts  [][]int
	rights [][]int

	counter int // for checking if a tree is projective

	n int
}

type depConsOpt func(*Dependency)

// FromAnnotatedSentence creates a dependency from an AnnotatedSentence.
func FromAnnotatedSentence(s AnnotatedSentence) depConsOpt {
	fn := func(d *Dependency) {
		wc := len(s)
		d.AnnotatedSentence = s
		d.wordCount = wc
		d.n = wc - 1
	}
	return fn
}

// AllocTree allocates the lefts and rights. Typical construction of the *Dependency doesn't allocate the trees as they're not necessary for a number of tasks.
func AllocTree() depConsOpt {
	fn := func(d *Dependency) {
		d.lefts = make([][]int, d.wordCount)
		d.rights = make([][]int, d.wordCount)
		for i := 0; i < d.wordCount; i++ {
			d.lefts[i] = make([]int, 0)
			d.rights[i] = make([]int, 0)
		}
	}
	return fn
}

// NewDependency creates a new *Dependency. It takes optional construction options:
//		FromAnnotatedSentence
//		AllocTree
func NewDependency(opts ...depConsOpt) *Dependency {
	d := new(Dependency)

	for _, opt := range opts {
		opt(d)
	}
	return d
}

func (d *Dependency) Sentence() AnnotatedSentence { return d.AnnotatedSentence }
func (d *Dependency) Lefts() [][]int              { return d.lefts }
func (d *Dependency) Rights() [][]int             { return d.rights }
func (d *Dependency) WordCount() int              { return d.wordCount }
func (d *Dependency) N() int                      { return d.n }

// please only use these for testing
func (d *Dependency) SetLefts(l [][]int)  { d.lefts = l }
func (d *Dependency) SetRights(r [][]int) { d.rights = r }

func (d *Dependency) Head(i int) int {
	if i < 0 || i >= d.wordCount || d.AnnotatedSentence[i].Head == nil {
		return -1
	}

	return d.AnnotatedSentence[i].HeadID()
}

func (d *Dependency) Label(i int) DependencyType {
	if i < 0 || i >= d.wordCount {
		return NoDepType
	}

	return d.AnnotatedSentence[i].DependencyType
}

func (d *Dependency) Annotation(i int) *Annotation {
	if i < 0 || i >= d.wordCount {
		return nullAnnotation
	}

	return d.AnnotatedSentence[i]
}

func (d *Dependency) AddArc(head, child int, label DependencyType) {
	d.AddChild(head, child)
	d.AddRel(child, label)
}

func (d *Dependency) AddChild(head, child int) {
	headAnn := d.AnnotatedSentence[head]
	d.AnnotatedSentence[child].SetHead(headAnn)

	if child < head {
		d.lefts[head] = append(d.lefts[head], child)
	} else {
		d.rights[head] = append(d.rights[head], child)
	}

	d.n++
}

func (d *Dependency) AddRel(child int, rel DependencyType) {
	// d.labels[child] = rel
	d.AnnotatedSentence[child].DependencyType = rel
}

func (d *Dependency) HasSingleRoot() bool {
	roots := 0
	for _, a := range d.AnnotatedSentence {
		h := a.HeadID()
		if h == 0 {
			roots++
		}
	}

	return roots == 1
}

func (d *Dependency) IsLegal() bool {
	var heads []int
	for _, a := range d.AnnotatedSentence {
		h := a.HeadID()
		if h < 0 || h > d.wordCount {
			return false
		}
		heads = append(heads, -1)
	}

	for i := 1; i < d.wordCount; i++ {
		for k := i; k > 0; {
			if heads[k] >= 0 && heads[k] < 1 {
				break
			}
			if heads[k] == i {
				return false
			}
			heads[k] = i
			k = d.AnnotatedSentence[k].HeadID()
		}
	}

	return true
}

func (d *Dependency) IsProjective() bool {
	d.counter = -1
	return d.projectiveVisit(0)
}

func (d *Dependency) projectiveVisit(w int) bool {
	for i := 1; i < w; i++ {
		if d.AnnotatedSentence[i].HeadID() == w && d.projectiveVisit(i) == false {
			return false
		}
	}

	d.counter++

	if w != d.counter {
		return false
	}

	for i := w + 1; i < d.wordCount; i++ {
		if d.AnnotatedSentence[i].HeadID() == w && d.projectiveVisit(i) == false {
			return false
		}
	}

	return true
}

func (d *Dependency) Root() int {
	for i := 1; i <= d.n; i++ {
		if d.Head(i) == 0 {
			return i
		}
	}

	return 0
}

func (d *Dependency) SprintRel() string {
	var buf bytes.Buffer

	for _, e := range d.Edges() {
		fmt.Fprintf(&buf, "%v(%q-%d, %q-%d)\n", e.Rel, e.Gov.Value, e.Gov.ID, e.Dep.Value, e.Dep.ID)
	}

	return buf.String()
}

type DependencyEdge struct {
	Gov *Annotation
	Dep *Annotation
	Rel DependencyType
}

// Sort interface

type edgeByID []DependencyEdge

func (b edgeByID) Len() int           { return len(b) }
func (b edgeByID) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b edgeByID) Less(i, j int) bool { return b[i].Gov.ID < b[j].Gov.ID }
