package lingo

import (
	"github.com/awalterschulze/gographviz"

	"fmt"

	"sync"
)

// A DependencyTree is an alternate form of representing a dependency parse.
// This form makes it easier to traverse the tree
type DependencyTree struct {
	Parent *DependencyTree

	ID   int            // the word number in a sentence
	Type DependencyType // refers to the dependency type to the parent
	Word *Annotation

	Children []*DependencyTree
}

func NewDependencyTree(parent *DependencyTree, ID int, ann *Annotation) *DependencyTree {
	return &DependencyTree{
		Parent:   parent,
		ID:       ID,
		Word:     ann,
		Children: make([]*DependencyTree, 0),
	}
}

func (d *DependencyTree) AddChild(child *DependencyTree) {
	d.Children = append(d.Children, child)
}

func (d *DependencyTree) AddRel(rel DependencyType) {
	d.Type = rel
}

func (d *DependencyTree) walk(c chan *DependencyTree, wg *sync.WaitGroup) {
	defer wg.Done()

	for _, child := range d.Children {
		wg.Add(1)
		go child.walk(c, wg)
	}
	c <- d // man someone should do somehting about my bad naming
}

func (d *DependencyTree) Dot() string {
	// walk graph
	c := make(chan *DependencyTree)
	out := make(chan string)

	go dotString(c, out)
	var wg sync.WaitGroup
	wg.Add(1)
	go d.walk(c, &wg)

	wg.Wait()
	close(c)
	return <-out
}

func dotString(c chan *DependencyTree, out chan string) {
	g := gographviz.NewEscape()
	g.SetName("G")
	g.SetDir(true) // it's always going to be a directed graph
	// g.AddNode("G", "Node_0x0", nil) // add the root

	for t := range c {
		id := fmt.Sprintf("Node_%p", t)
		attrs := map[string]string{
			"label": fmt.Sprintf("%d: \"%s/%s\"", t.ID, t.Word.Value, t.Word.POSTag),
		}
		g.AddNode("G", id, attrs)

		if t.Parent == nil {
			continue
		}

		parentID := fmt.Sprintf("Node_%p", t.Parent)
		edgeAttrs := map[string]string{
			"label": fmt.Sprintf("%v", t.Type),
		}
		g.AddEdge(parentID, id, true, edgeAttrs)
	}
	out <- g.String()
}

func (d *DependencyTree) Walk(fn func(interface{})) {
	for _, child := range d.Children {
		child.Walk(fn)
	}

	if fn != nil {
		fn(d)
	}
}
