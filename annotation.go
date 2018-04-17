package lingo

import (
	"errors"
	"fmt"
	"strings"
)

// Annotation is the word and it's metadata.
// This includes the position, its dependency head (if available), its lemma, POSTag, etc
//
// A collection of Annoations - AnnotatedSentence is also a representation of a dependency parse
//
// Every field is exported for easy gobbing. be very careful with setting stuff
type Annotation struct {
	Lexeme
	POSTag
	// NER

	// fields to do with an annotation being in a collection
	DependencyType
	ID       int
	Head     *Annotation
	children AnnotationSet //will not be serialized

	// info about the annotation itself
	Lemma   string
	Lowered string
	Stem    string

	// auxiliary data for processing
	Cluster
	Shape
	WordFlag
}

func NewAnnotation() *Annotation {
	return &Annotation{
		Lexeme: nullLexeme,
		Lemma:  "",
		Shape:  Shape(""),
	}
}

// AnnotationFromLexTag is only ever used in tests. Fixer is optional
func AnnotationFromLexTag(l Lexeme, t POSTag, f AnnotationFixer) *Annotation {
	a := &Annotation{
		Lexeme:         l,
		POSTag:         t,
		DependencyType: NoDepType,
		Lemma:          "",
		Lowered:        strings.ToLower(l.Value),
	}

	// it's ok to panic - it will cause the tests to fail
	if err := a.Process(f); err != nil {
		panic(err)
	}

	return a
}

func (a *Annotation) Clone() *Annotation {
	b := *a
	b.ID = -1
	b.Head = nil
	b.children = nil
	b.DependencyType = NoDepType

	return &b
}

func (a *Annotation) SetHead(headAnn *Annotation) {
	a.Head = headAnn
	if headAnn != rootAnnotation && headAnn != startAnnotation && headAnn != nullAnnotation {
		headAnn.children = append(headAnn.children, a)
	}
}

func (a *Annotation) HeadID() int {
	if a.Head != nil {
		return a.Head.ID
	}
	return -1
}

func (a *Annotation) IsNumber() bool {
	return IsNumber(a.POSTag) && (a.LexemeType != Date && a.LexemeType != Time && a.LexemeType != URI)
}

func (a *Annotation) String() string {
	return a.Value
}

func (a *Annotation) GoString() string {
	s := fmt.Sprintf("%q/%s", a.Lexeme.Value, a.POSTag)

	if a.Head != nil {
		return fmt.Sprintf("(%v) <-%v- (%q/%s) ", s, a.DependencyType, a.Head.Value, a.Head.POSTag)
	}
	return s
}

func (a *Annotation) Process(f AnnotationFixer) error {
	if a.Lexeme != nullLexeme {
		a.Lowered = strings.ToLower(a.Value)
		a.Shape = a.Lexeme.Shape()
		a.WordFlag = a.Lexeme.Flags()

		var err error
		if f != nil {
			var stem string
			if stem, err = f.Stem(a.Lowered); err != nil {
				if _, ok := err.(componentUnavailable); !ok {
					return err
				}
			}
			a.Stem = stem

			var clust map[string]Cluster
			if clust, err = f.Clusters(); err == nil {
				a.Cluster = clust[a.Value]
			}
		}

		return nil
	}
	return errors.New("No Lexeme!")
}

var rootAnnotation = &Annotation{
	Lexeme:         rootLexeme,
	POSTag:         ROOT_TAG,
	DependencyType: Root,
	ID:             0,
	Head:           nil,
	Lemma:          "",
	Lowered:        "",
	Cluster:        0,
	Shape:          "",
	WordFlag:       NoFlag,
}

var startAnnotation = &Annotation{
	Lexeme:         startLexeme,
	POSTag:         ROOT_TAG,
	DependencyType: NoDepType,
	ID:             -1,
	Head:           nil,
	Lemma:          "",
	Lowered:        "",
	Cluster:        0,
	Shape:          "",
	WordFlag:       NoFlag,
}

var nullAnnotation = &Annotation{
	Lexeme:         nullLexeme,
	POSTag:         X,
	DependencyType: NoDepType,
	ID:             -1,
	Head:           nil,
	Lemma:          "",
	Lowered:        "",
	Cluster:        0,
	Shape:          "",
	WordFlag:       NoFlag,
}

func RootAnnotation() *Annotation  { return rootAnnotation }
func StartAnnotation() *Annotation { return startAnnotation }
func NullAnnotation() *Annotation  { return nullAnnotation }

func StringToAnnotation(s string, f AnnotationFixer) *Annotation {
	l := MakeLexeme(s, Word)
	a := NewAnnotation()
	a.Lexeme = l
	if err := a.Process(f); err != nil {
		panic(err.Error())
	}
	return a
}

type AnnotationFixer interface {
	Lemmatizer
	Stemmer
	Clusters() (map[string]Cluster, error)
}
