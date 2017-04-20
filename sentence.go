package lingo

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/pkg/errors"
)

/* Lexeme Sentence */
type LexemeSentence []Lexeme

func NewLexemeSentence() LexemeSentence { return LexemeSentence(make([]Lexeme, 0)) }

func (ls LexemeSentence) String() string {
	var buf bytes.Buffer
	for _, lex := range ls {
		buf.WriteString(lex.Value)
		buf.WriteString(" ")
	}
	return strings.Trim(buf.String(), " ")
}

/* Annotated Sentence */

// AnnotatedSentence is a sentence, but each word has been annotated.
type AnnotatedSentence []*Annotation

func NewAnnotatedSentence() AnnotatedSentence { return make(AnnotatedSentence, 0) }

func (as AnnotatedSentence) Clone() AnnotatedSentence {
	retVal := make(AnnotatedSentence, len(as))

	for i, a := range as {
		// don't clone rootAnnotation
		if i == 0 && a == rootAnnotation {
			retVal[i] = a
			continue
		}
		retVal[i] = a.Clone()
	}
	return retVal
}

func (as AnnotatedSentence) SetID() {
	for i, a := range as {
		if i == 0 && a == rootAnnotation {
			continue
		}
		a.ID = i
	}
}

func (as AnnotatedSentence) Fix() {
	if as[0].Lexeme == rootLexeme {
		as[0] = rootAnnotation
	}

	as.SetID()

	for _, a := range as {
		if a.Head != nil {
			if a.HeadID() == -1 && a.Head.Lexeme == rootLexeme {
				a.Head = rootAnnotation
				continue
			}
			a.SetHead(as[a.HeadID()])
		}
	}
}

func (as AnnotatedSentence) IsValid() bool {
	// check that IDs are set
	zeroes := 0
	for _, a := range as {
		if a.ID == 0 {
			zeroes++
		}
	}
	// IDs not properly set
	if zeroes > 1 {
		return false
	}

	// TODO
	// check that there is only one root

	return true
}

/* Return slices of x */

// Phrase returns the slice of the sentence. While you can do the same by simply doing as[start:end], this method returns errors instead of panicking
func (as AnnotatedSentence) Phrase(start, end int) (AnnotatedSentence, error) {
	if start < 0 {
		return nil, errors.Errorf("Start: %d < 0", start)
	}
	if end > len(as) {
		return nil, errors.Errorf("End: %d > len(as): %d", end, len(as))
	}
	return as[start:end], nil
}

// IDs returns the list of IDs in the sentence. The return value has exactly the same length as the sentence.
func (as AnnotatedSentence) IDs() []int {
	retVal := make([]int, len(as))
	for i, a := range as {
		retVal[i] = a.ID
	}
	return retVal
}

// Tags returns the POSTags of the sentence. The return value has exactly the same length as the sentence.
func (as AnnotatedSentence) Tags() []POSTag {
	retVal := make([]POSTag, len(as))
	for i, a := range as {
		retVal[i] = a.POSTag
	}
	return retVal
}

// Heads returns the head IDs of the sentence. The return value has exactly the same length as the sentence.
func (as AnnotatedSentence) Heads() []int {
	retVal := make([]int, len(as))
	for i, a := range as {
		retVal[i] = a.HeadID()
	}
	return retVal
}

// Leaves returns the *Annotations which are leaves. If the dependency hasn't been set yet, every single *Annotation is a leaf.
func (as AnnotatedSentence) Leaves() (retVal []int) {
	for i := range as {
		if len(as.Children(i)) == 0 {
			retVal = append(retVal, i)
		}
	}
	return
}

// Labels returns the DependencyTypes of the sentence. The return value has exactly the same length as the sentence.
func (as AnnotatedSentence) Labels() []DependencyType {
	retVal := make([]DependencyType, len(as))
	for i, a := range as {
		retVal[i] = a.DependencyType
	}
	return retVal
}

// StringSlice returns the original words as a slice of string. The return value has exactly the same length as the sentence.
func (as AnnotatedSentence) StringSlice() []string {
	retVal := make([]string, len(as), len(as))
	for i, a := range as {
		retVal[i] = a.Value
	}
	return retVal
}

// LoweredStringSlice returns the lowercased version of the words in the sentence as a slice of string. The return value has exactly the same length as the sentence.
func (as AnnotatedSentence) LoweredStringSlice() []string {
	retVal := make([]string, len(as), len(as))
	for i, a := range as {
		retVal[i] = a.Lowered
	}
	return retVal
}

// Lemmas returns the lemmas as as slice of string. The return value has exactly the same length as the sentence.
func (as AnnotatedSentence) Lemmas() []string {
	lemmas := make([]string, len(as))
	for i, a := range as {
		lemmas[i] = a.Lemma
	}
	return lemmas
}

// Stems returns the stems as a slice of string. The return value has exactly the same length as the sentence.
func (as AnnotatedSentence) Stems() []string {
	stems := make([]string, len(as))
	for i, a := range as {
		stems[i] = a.Stem
	}
	return stems
}

func (as AnnotatedSentence) Children(h int) (retVal []int) {
	for i, v := range as {
		if v.HeadID() == h {
			retVal = append(retVal, i)
		}
	}
	return
}

func (as AnnotatedSentence) Edges() (retVal []DependencyEdge) {
	for _, a := range as {
		var head = -1

		if a.Head != nil {
			head = a.HeadID()
		}

		if head == -1 {
			head = 0
		}
		edge := DependencyEdge{as[head], a, a.DependencyType}
		retVal = append(retVal, edge)
	}
	sort.Sort(edgeByID(retVal))
	return
}

/* To other structures */

func (as AnnotatedSentence) Dependency() *Dependency {
	return NewDependency(FromAnnotatedSentence(as))
}

func (as AnnotatedSentence) Tree() *DependencyTree {
	tracker := make([]*DependencyTree, len(as))

	rootNode := NewDependencyTree(nil, 0, rootAnnotation)
	tracker[0] = rootNode

	for i := 1; i < len(as); i++ {
		head := as[i].HeadID()
		var headDep *DependencyTree

		if head == -1 {
			headDep = rootNode
		} else {
			headDep = tracker[head]
		}

		if headDep == nil {
			// make a dependency for the head
			headDep = NewDependencyTree(nil, head, as[head])
			tracker[head] = headDep
		}

		dep := tracker[i]

		if dep == nil {
			dep = NewDependencyTree(headDep, i, as[i])
			tracker[i] = dep
		} else {
			dep.Parent = headDep
		}

		headDep.AddChild(dep)
		dep.Type = as[i].DependencyType

	}
	// return tracker[len(tracker)-1]
	// log.Printf("Tracker: %v, len(as): %d. Root: %v", tracker, len(as), rootNode.Children)
	return rootNode
}

// Stringer interface

func (as AnnotatedSentence) String() string {
	var buf bytes.Buffer
	for i, a := range as {
		buf.WriteString(fmt.Sprintf("%s/%s", a.Value, a.POSTag))
		if i < len(as)-1 {
			buf.WriteString(" ")
		}
	}
	return buf.String()
}

func (as AnnotatedSentence) ValueString() string {
	var buf bytes.Buffer
	for i, a := range as {
		buf.WriteString(a.Value)
		if i < len(as)-1 {
			buf.WriteString(" ")
		}
	}
	return buf.String()
}

func (as AnnotatedSentence) LoweredString() string {
	var buf bytes.Buffer
	for i, a := range as {
		buf.WriteString(a.Lowered)
		if i < len(as)-1 {
			buf.WriteString(" ")
		}
	}
	return buf.String()
}

func (as AnnotatedSentence) LemmaString() string {
	var buf bytes.Buffer
	for i, a := range as {
		buf.WriteString(a.Lemma)
		if i < len(as)-1 {
			buf.WriteString(" ")
		}
	}
	return buf.String()
}

func (as AnnotatedSentence) StemString() string {
	var buf bytes.Buffer
	for i, a := range as {
		buf.WriteString(a.Stem)
		if i < len(as)-1 {
			buf.WriteString(" ")
		}
	}
	return buf.String()
}

// sort interface
func (as AnnotatedSentence) Len() int           { return len(as) }
func (as AnnotatedSentence) Swap(i, j int)      { as[i], as[j] = as[j], as[i] }
func (as AnnotatedSentence) Less(i, j int) bool { return as[i].ID < as[j].ID }
