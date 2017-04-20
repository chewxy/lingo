package lingo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

type dummyAnnotation struct {
	POSTag         `json:"POSTag"`
	DependencyType `json:"Label"`

	ID    int    `json:"ID"`
	Head  int    `json:"Head"`
	Value string `json:"Value"`
	Lemma string `json:"Lemma"`
	Stem  string `json:"Stem"`

	Cluster  `json:"Cluster"`
	Shape    `json:"Shape"`
	WordFlag `json:"WordFlat"`
}

// func (a *Annotation) MarshalText() ([]byte, error) {
// 	var buf bytes.Buffer
// 	if a.Head != nil {
// 		fmt.Fprintf(&buf, "%v(%q/%v-%d, %q/%v-%d)", a.DependencyType, a.Value, a.POSTag, a.ID, a.Head.Value, a.Head.POSTag, a.Head.ID)
// 	} else if a == rootAnnotation {
// 		fmt.Fprintf(&buf, "ROOT")
// 	} else {
// 		fmt.Fprintf(&buf, "%q/%v-%d", a.Value, a.POSTag, a.ID)
// 	}
// 	return buf.Bytes(), nil
// }

func (a *Annotation) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteRune('{')

	fmt.Fprintf(&buf, "\"ID\": %d,", a.ID)
	fmt.Fprintf(&buf, "\"Value\": %q,", a.Value)
	fmt.Fprintf(&buf, "\"POSTag\": \"%v\",", a.POSTag)
	fmt.Fprintf(&buf, "\"Label\": \"%v\"", a.DependencyType)

	if a.Head != nil {
		if a.Head == rootAnnotation {
			fmt.Fprintf(&buf, ", \"Head\": -1000") // special signifier for root annotations
		} else {
			fmt.Fprintf(&buf, ", \"Head\": %d", a.HeadID())
		}
	}

	if a.Lemma != "" {
		fmt.Fprintf(&buf, ", \"Lemma\": %q", a.Lemma)
	}

	// Lowered is not serialized because it's a simple function call away

	if a.Stem != "" {
		fmt.Fprintf(&buf, ",\"Stem\": %q", a.Stem)
	}

	if a.Cluster > 0 {
		fmt.Fprintf(&buf, ",\"Cluster\": %d", a.Cluster)
	}

	if a.Shape != "" {
		fmt.Fprintf(&buf, ",\"Shape\": %q", a.Shape)
	}

	if a.WordFlag > 0 {
		fmt.Fprintf(&buf, ",\"WordFlag\": %d", a.WordFlag)
	}
	buf.WriteRune('}')
	return buf.Bytes(), nil
}

func (a *Annotation) UnmarshalJSON(b []byte) error {
	if a == nil {
		// error
		return errors.Errorf("Cannot unmarshal json to a nul")
	}

	d := dummyAnnotation{}
	if err := json.Unmarshal(b, &d); err != nil {
		return err
	}

	a.Value = d.Value
	a.POSTag = d.POSTag
	a.DependencyType = d.DependencyType
	a.ID = d.ID
	a.Lemma = d.Lemma
	a.Stem = d.Stem
	a.Cluster = d.Cluster
	a.Shape = d.Shape
	a.WordFlag = d.WordFlag

	return nil
}

func (as AnnotatedSentence) MarshalJSON() ([]byte, error) {
	buf := new(bytes.Buffer)
	encoder := json.NewEncoder(buf)

	buf.WriteRune('[')
	for i, a := range as {
		if err := encoder.Encode(a); err != nil {
			return nil, err
		}
		if i < len(as)-1 {
			buf.WriteRune(',')
		}
	}
	buf.WriteRune(']')
	return buf.Bytes(), nil
}

func (as *AnnotatedSentence) UnmarshalJSON(b []byte) error {
	dummies := make([]dummyAnnotation, 0)

	if err := json.Unmarshal(b, &dummies); err != nil {
		return err
	}

	asL := len(*as)
	l := len(dummies)
	if asL != l {
		diff := l - asL
		(*as) = append(*as, make(AnnotatedSentence, diff)...)
	}

	for i, d := range dummies {
		a := (*as)[i]
		if d.Value == "-ROOT-" {
			(*as)[i] = rootAnnotation
			continue
		}

		if a == nil {
			a = new(Annotation)
		}

		a.Value = d.Value
		a.POSTag = d.POSTag
		a.DependencyType = d.DependencyType
		a.ID = d.ID
		a.Lemma = d.Lemma
		a.Stem = d.Stem
		a.Cluster = d.Cluster
		a.Shape = d.Shape
		a.WordFlag = d.WordFlag

		(*as)[i] = a
	}

	// fix up head IDs
	for i, d := range dummies {
		a := (*as)[i]
		head := d.Head
		if head == -1000 {
			a.SetHead(rootAnnotation)
		} else {
			a.SetHead((*as)[head])
		}
	}

	// TODO: fix up other things
	for _, a := range *as {
		a.Lowered = strings.ToLower(a.Value)
	}

	return nil
}
