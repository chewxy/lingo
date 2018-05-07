package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"

	"github.com/chewxy/lingo"
	"github.com/chewxy/lingo/dep"
	"github.com/chewxy/lingo/pos"
)

var save = flag.String("save", "", "save as...")
var load = flag.String("load", "", "load a model")
var loadPOS = flag.String("PTmodel", "", "load a POS Tagger model")
var clusterFiles = flag.String("cluster", "", "Brown Cluster files. If nothing is passed in, then the brown cluster won't be used")
var trainFile = flag.String("train", "", "Training on... (Only CONLLU formatted training files are accepted)")
var testFile = flag.String("test", "", "Test on... (Only CONLLU formatted training files are accepted). If this is not provided, the model will be trained without crossvalidation")
var cv = flag.Bool("cv", false, "Cross validate training model? Defaults to false.")
var epoch = flag.Int("epoch", 10, "Training epochs. Defaults to 10")
var format = flag.String("f", "", "Format to output. Default is none. Accepts: {json, dot}")

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var memprofile = flag.String("memprofile", "", "write memory profile to this file")

var clusters map[string]lingo.Cluster
var POSModel *pos.Model
var DepModel *dep.Model
var toLoad, toTrain bool

func init() {
	if lingo.BUILD_TAGSET != "stanfordtags" && lingo.BUILD_TAGSET != "universaltags" {
		log.Fatalf("Tagset %q unsupported", lingo.BUILD_TAGSET)
	}

	if lingo.BUILD_RELSET != "stanfordrel" && lingo.BUILD_RELSET != "universalrel" {
		log.Fatalf("Relset %q unsupported", lingo.BUILD_RELSET)
	}
}

func cleanup(sigChan chan os.Signal, cpuprofiling, memprofiling bool) {
	select {
	case <-sigChan:
		log.Println("EMERGENCY EXIT")
		if cpuprofiling {
			pprof.StopCPUProfile()

		}
		if memprofiling {
			f, err := os.Create(*memprofile)
			if err != nil {
				log.Fatal(err)
			}
			pprof.WriteHeapProfile(f)
			f.Close()
		}
		saveModel()
		os.Exit(1)
	}
}

func main() {
	flag.Parse()
	validateFlags()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	var cpuprofiling, memprofiling bool
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		cpuprofiling = true
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if *memprofile != "" {
		memprofiling = true
	}

	go cleanup(sigChan, cpuprofiling, memprofiling)

	loadPOSModel()
	if toLoad {
		loadDepModel()
	}

	if toTrain {
		loadTreebanks()
		train()
	}

	saveModel()
}
