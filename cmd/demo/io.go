package main

import (
	"log"
	"os"

	"github.com/chewxy/lingo"
	"github.com/chewxy/lingo/dep"
	"github.com/chewxy/lingo/pos"
)

const (
	posModelFile = `model/pos_stanfordtags_universalrel.final.model`
	depModelFile = `model/dep_stanfordtags_universalrel.final.model`
	brownCluster = `clusters.txt`
)

func io() {
	var err error
	log.Println("loading POS Tagger model")
	if posModel, err = pos.Load(posModelFile); err != nil {
		log.Fatal(err)
	}

	log.Println("loading Dependency Parser model")
	if depModel, err = dep.Load(depModelFile); err != nil {
		log.Fatal(err)
	}
	var f *os.File
	if f, err = os.Open(brownCluster); err != nil {
		log.Fatal(err)
	}
	clusters = lingo.ReadCluster(f)
}
