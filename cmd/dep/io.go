package main

import (
	"log"

	"github.com/chewxy/lingo/dep"
	"github.com/chewxy/lingo/pos"
	"github.com/chewxy/lingo/treebank"
)

func validateFlags() {
	if *load == "" && *trainFile == "" {
		log.Fatal("Must either load a model or pass in a training file")
	}

	if *epoch < 0 {
		log.Fatal("epochs must only be positive numbers")
	}

	if *load != "" {
		toLoad = true
	}

	if *trainFile != "" {
		toTrain = true
	}

	if *testFile != "" {
		*cv = true
	}

	// warnings
	if *load == "" && *save == "" {
		log.Println("WARNING: Models that have been trained will NOT be saved")
	}
}

func loadTreebanks() {
	if *trainFile != "" {
		trainTB = treebank.LoadUniversal(*trainFile)
	}

	if *testFile != "" {
		testTB = treebank.LoadUniversal(*testFile)
	}
}

func loadPOSModel() {
	var err error
	if *loadPOS == "" {
		log.Fatal("Cannot proceed without having a POS model")
	}
	if POSModel, err = pos.Load(*loadPOS); err != nil {
		log.Fatal(err)
	}
}

func loadDepModel() {
	var err error

	if DepModel, err = dep.Load(*load); err != nil {
		log.Fatal(err)
	}
}

func saveModel() {
	if *save != "" && DepModel != nil {
		DepModel.Save(*save)
	}
}
