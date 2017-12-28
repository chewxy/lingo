package main

import (
	"log"

	"github.com/chewxy/lingo/dep"
	"github.com/chewxy/lingo/treebank"
	"gorgonia.org/tensor"
)

var trainTB []treebank.SentenceTag
var testTB []treebank.SentenceTag

func train() {
	conf := dep.DefaultNNConfig
	conf.Dtype = tensor.Float32
	var trainer *dep.Trainer

	if testTB != nil {
		log.Printf("TRAINING WITH CROSSVALIDATION")
		trainer = dep.NewTrainer(dep.WithGeneratedCorpus(trainTB...), dep.WithTrainingSet(trainTB), dep.WithCrossValidationSet(testTB), dep.WithConfig(conf))
		trainer.SaveBest = "TMP.model"
		if err := trainer.Init(); err != nil {
			log.Fatalf("Unable to initialize trainer: \n%+v", err)
		}

		prog := trainer.Perf()
		cost := trainer.Cost()
		go func() {
			for {
				select {
				case p := <-prog:
					log.Printf("%v\n", p)
				case c := <-cost:
					log.Printf("Cost %v\n", c)
				}
			}
		}()

	} else {
		trainer = dep.NewTrainer(dep.WithGeneratedCorpus(trainTB...), dep.WithTrainingSet(trainTB), dep.WithConfig(conf))
		if err := trainer.Init(); err != nil {
			log.Fatalf("Unable to initialize trainer: \n%+v", err)
		}

		prog := trainer.Cost()
		go func() {
			for cost := range prog {
				log.Printf("Cost %v\n", cost)
			}
		}()
	}

	if err := trainer.Train(*epoch); err != nil {
		log.Fatal(err)
	}

	DepModel = trainer.Model
}
