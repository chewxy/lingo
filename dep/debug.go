// +build debug

package dep

import (
	"bytes"
	"fmt"
	"log"
	"runtime"
	"strings"
	"sync/atomic"

	"github.com/chewxy/lingo"
)

const BUILD_DEBUG = "PARSER: DEBUG BUILD"
const BUILD_DIAG = "Diagnostic Build"

const DEBUG = true

var READMEMSTATS = true

var TABCOUNT uint32 = 0

func tabcount() int {
	return int(atomic.LoadUint32(&TABCOUNT))
}

func enterLoggingContext() {
	atomic.AddUint32(&TABCOUNT, 1)
	tc := tabcount()
	log.SetPrefix(strings.Repeat("\t", tc))
}

func leaveLoggingContext() {
	tc := tabcount()
	tc--

	if tc < 0 {
		atomic.StoreUint32(&TABCOUNT, 0)
		tc = 0
	} else {
		atomic.StoreUint32(&TABCOUNT, uint32(tc))
	}
	log.SetPrefix(strings.Repeat("\t", tc))
}

func logf(format string, others ...interface{}) {
	if !DEBUG {
		return
	}
	log.Printf(format, others...)
}

func logTrainingProgress(iteration, correct, total, length, possibles int) {
	if !DEBUG {
		return
	}

	log.Printf("Iteration %d. Correct/Total: %d/%d = %.2f", iteration, correct, total, float64(correct)/float64(total))
	log.Printf("DictSize: %d/%d, load factor of: %.2f", length, possibles, float64(length)/float64(possibles))
}

func logMemStats() {
	if !DEBUG || !READMEMSTATS {
		return
	}

	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	log.Printf("Allocated          : %.2f MB", (float64(mem.Alloc)/1024)/float64(1024))
	log.Printf("Total Allocated    : %.2f MB", (float64(mem.TotalAlloc)/1024)/float64(1024))
	log.Printf("Heap Allocted      : %.2f MB", (float64(mem.HeapAlloc)/1024)/float64(1024))
	log.Printf("Sys Total Allocated: %.2f MB", (float64(mem.HeapSys)/1024)/float64(1024))
	log.Println("----------")
}

func recoverFrom(format string, attrs ...interface{}) {
	if r := recover(); r != nil {
		log.Printf(format, attrs...)
		panic(r)
	}
}

/* Nice output of shit */
func (d *Parser) SprintFeatures(features []int) string {
	// tabcount := int(atomic.LoadUint32(&TABCOUNT))

	var buf bytes.Buffer

	for i := 0; i < 18; i++ {
		number := features[i]
		id := number - wordFeatsStartAt
		word, _ := d.corpus.Word(id)

		if word == "" {
			word = "-NULL-"
		}

		buf.WriteString(fmt.Sprintf("%d, %q, %d \n", feature(i), word, number))
	}

	for i := 0; i < 18; i++ {
		number := features[i+18]

		buf.WriteString(fmt.Sprintf("%d, %v, %d\n", feature(i+18), lingo.POSTag(number), number))
	}

	for i := 0; i < 12; i++ {
		number := features[i+36]
		id := number - labelFeatsStartAt

		buf.WriteString(fmt.Sprintf("%d, %v, %d\n", feature(i+36), lingo.DependencyType(id), number))
	}

	return buf.String()
}

func SprintScores(scores []float64, ts []transition) string {
	var buf bytes.Buffer
	for i, v := range scores {
		if i >= len(ts) {
			buf.WriteString(fmt.Sprintf("UNKNOWN TRANSITION, %v\n", v))
			continue
		}
		buf.WriteString(fmt.Sprintf("%v, %v\n", ts[i], v))
	}
	return buf.String()
}

func SprintFloatSlice(a []float64) string {
	var buf bytes.Buffer
	buf.WriteString("[")
	for i, v := range a {
		if i < len(a)-1 {
			buf.WriteString(fmt.Sprintf("%v, ", v))
		} else {
			buf.WriteString(fmt.Sprintf("%v", v))
		}
	}
	buf.WriteString("]")
	return buf.String()
}
