package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime/pprof"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/chewxy/lingo"
	"github.com/chewxy/lingo/lexer"
	"github.com/chewxy/lingo/pos"
	"github.com/chewxy/lingo/treebank"
)

var save = flag.String("save", "", "save as...")
var load = flag.String("load", "", "load a model")
var clusterFiles = flag.String("cluster", "", "Brown Cluster files. If nothing is passed in, then the brown cluster won't be used")
var trainFile = flag.String("train", "", "Training on... files that end with '.conllu' will be treated as CONLLU formatted files. Files ending with '.zip' will be treted as EWT files")
var testFile = flag.String("test", "", "Test on... Files to cross validate the model on. If this is provided, automatic crossvalidation will be done")
var cv = flag.Bool("cv", false, "Cross validate training model? Defaults to false.")
var epoch = flag.Int("epoch", 1500, "Training epochs. Defaults to 1500")
var inspect = flag.String("inpect", "", "Inspect all the wrong outputs to figure out what went wrong in the POSTagging. This is useful for debugging")
var input = flag.String("input", "", "Input sentence to tag")

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var memprofile = flag.String("memprofile", "", "write memory profile to this file")

var clusters map[string]lingo.Cluster
var model *pos.Model

func receive(sentences chan lingo.AnnotatedSentence, wg *sync.WaitGroup) {
	defer wg.Done()
	for sent := range sentences {
		for _, a := range sent {
			fmt.Printf("%#v: %s| %s | %s | %d\n", a, a.POSTag, a.Lemma, a.WordFlag, a.Cluster)
		}
	}
}

func pipeline(s string) {
	l := lexer.New(s, strings.NewReader(s))
	pt := pos.New(pos.WithModel(model))

	pt.Input = l.Output
	var wg sync.WaitGroup

	go l.Run()
	go receive(pt.Output, &wg)

	wg.Add(1)

	pt.Run()
	wg.Wait()
}

func validateFlags() {
	if *load == "" && *trainFile == "" {
		log.Fatal("Must either load a model or pass in a training file")
	}

	if *epoch < 0 {
		log.Fatal("epochs must be positive numbers only!")
	}

	if *testFile != "" {
		*cv = true
	}

	// warnings

	if *load == "" && *save == "" {
		log.Println("WARNING: Models that are trained will NOT be saved")
	}
}

func loadOrTrain() {
	var trained *pos.Tagger
	if *clusterFiles != "" {
		f, err := os.Open(*clusterFiles)
		if err != nil {
			log.Fatal(err)
		}
		clusters = lingo.ReadCluster(f)

		trained = pos.New(pos.WithCluster(clusters), pos.WithStemmer(stemmer{}))
	} else {
		trained = pos.New()
	}

	if *load != "" {
		start := time.Now()
		var err error
		if model, err = pos.Load(*load); err != nil {
			log.Fatal(err)
		}
		log.Printf("Loading model from %q took %v", *load, time.Since(start))
		return
	}

	var sentences []treebank.SentenceTag
	switch {
	case strings.HasSuffix(*trainFile, ".zip"):
		sentences = treebank.LoadEWT(*trainFile)

		// TODO split sentences for crossvalidation

	case strings.HasSuffix(*trainFile, ".conllu"):
		sentences = treebank.LoadUniversal(*trainFile)
	default:
		f, err := os.Open(*trainFile)
		if err != nil {
			log.Fatal(err)
		}

		sentences = treebank.ReadConllu(f)
	}

	log.Printf("Start training for %d epochs...", *epoch)
	start := time.Now()
	trained.Train(sentences, *epoch)
	log.Printf("End Training. Training took %v minutes", time.Since(start).Minutes())

	if *save != "" {
		trained.Save(*save)
		log.Printf("Model saved as: %v", *save)
	}
}

func cleanup(sigChan chan os.Signal, profiling bool) {
	select {
	case <-sigChan:
		log.Println("EMERGENCY EXIT")
		if profiling {
			pprof.StopCPUProfile()
		}
		os.Exit(1)
	}
}

func main() {
	flag.Parse()

	if lingo.BUILD_TAGSET != "stanfordtags" && lingo.BUILD_TAGSET != "universaltags" {
		log.Fatalf("Tagset: %v is unsupported", lingo.BUILD_TAGSET)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	var profiling bool
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		profiling = true
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	go cleanup(sigChan, profiling)

	validateFlags()
	loadOrTrain()

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.WriteHeapProfile(f)
		f.Close()
	}

	if *input != "" {
		pipeline(*input)
	}

	if *cv {
		log.Printf("Cross Validating now")
		testSentences := treebank.LoadUniversal(*testFile)
		testModel(testSentences)
	}

}
