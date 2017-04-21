package main

import (
	"fmt"
	"strings"

	"github.com/chewxy/lingo"
	"github.com/chewxy/lingo/dep"
	"github.com/chewxy/lingo/lexer"
	"github.com/chewxy/lingo/pos"
	"github.com/kljensen/snowball"
	"github.com/pkg/errors"
)

var posModel *pos.Model
var depModel *dep.Model

var clusters map[string]lingo.Cluster

type stemmer struct{}

func (stemmer) Stem(a string) (string, error) {
	return snowball.Stem(a, "english", true)
}

type fixer struct {
	stemmer
}

func (f fixer) Clusters() (map[string]lingo.Cluster, error) { return clusters, nil }
func (f fixer) Lemmatize(a string, pt lingo.POSTag) ([]string, error) {
	return nil, nocomp("lemmatizer")
}

type nocomp string

func (e nocomp) Error() string     { return fmt.Sprintf("no %v", string(e)) }
func (e nocomp) Component() string { return string(e) }

func pipeline(s string) (d *lingo.Dependency, err error) {
	if posModel == nil || depModel == nil {
		return nil, errors.Errorf("Unable to create a pipeline")
	}
	lx := lexer.New(s, strings.NewReader(s))
	pt := pos.New(pos.WithModel(posModel), pos.WithStemmer(stemmer{}))
	dp := dep.New(depModel)

	// pipeline
	pt.Input = lx.Output
	dp.Input = pt.Output

	go lx.Run()
	go pt.Run()
	go dp.Run()

	var ok bool
	for {
		select {
		case d, ok = <-dp.Output:
			if !ok {
				continue
			}
			return
		case err = <-dp.Error:
			return
		}
	}
}
