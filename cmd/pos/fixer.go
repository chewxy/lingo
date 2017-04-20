// +build !chewxy

package main

import (
	"fmt"

	"github.com/chewxy/lingo"
	"github.com/kljensen/snowball"
)

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
