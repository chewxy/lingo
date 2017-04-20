package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/chewxy/lingo"
	"github.com/chewxy/lingo/lexer"
	"github.com/chewxy/lingo/pos"
	"github.com/chewxy/lingo/treebank"
)

type testResult struct {
	tagged lingo.AnnotatedSentence
	actual lingo.AnnotatedSentence
}

func (tr testResult) compare() (int, bool) {
	tagged := tr.tagged
	actual := tr.actual

	var sameLength bool = true

	if len(tagged) != len(actual) {
		sameLength = false
	}

	var counter int
	for i, v := range actual {
		if i >= len(tagged) {
			break
		}
		if v.POSTag == tagged[i].POSTag {
			counter++
		}
	}
	return counter, sameLength
}

func crossValidate(resultChan chan testResult) {
	diffLengthCount := 0
	totalLength := 0
	correctCount := 0
	sentences := 0

	var wrongResults []testResult

	for res := range resultChan {
		sentences++
		length := len(res.actual)
		cc, sl := res.compare()
		if !sl {
			diffLengthCount++
		}
		correctCount += cc
		totalLength += length

		if cc != length && *inspect != "" {
			wrongResults = append(wrongResults, res)
		}
	}

	if *inspect != "" {
		f, err := os.OpenFile(*inspect, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Fatal(err)
		}

		// can write directly to f
		var buf bytes.Buffer
		for _, res := range wrongResults {
			fmt.Fprintf(&buf, "Sentence: \nW:%v\nG:%v\nTags:\nW: %v\nG: %v\n\n", res.actual.StringSlice(), res.tagged.StringSlice(), res.actual.Tags(), res.tagged.Tags())
		}

		f.WriteString(buf.String())
		f.Close()
	}

	fmt.Printf("CrossValidation: %d/%d = %f. Differing Lengths : %d/%d = %f\n", correctCount, totalLength, float64(correctCount)/float64(totalLength), diffLengthCount, sentences, float64(diffLengthCount)/float64(sentences))
}

func collect(ch chan lingo.AnnotatedSentence, correct lingo.AnnotatedSentence, outCh chan testResult, wg *sync.WaitGroup) {
	defer wg.Done()

	for sentence := range ch {
		outCh <- testResult{sentence, correct}
	}
}

func testModel(sentences []treebank.SentenceTag) {
	resultChan := make(chan testResult)

	go func() {
		defer close(resultChan)
		var wg sync.WaitGroup
		for _, sentence := range sentences {
			wg.Add(1)
			input := sentence.String()
			correct := sentence.AnnotatedSentence(fixer{stemmer{}})
			ch := make(chan lingo.AnnotatedSentence)
			go collect(ch, correct, resultChan, &wg)
			go cvpipeline(input, ch)
		}
		wg.Wait()
	}()

	crossValidate(resultChan)

}

func cvpipeline(s string, output chan lingo.AnnotatedSentence) {
	l := lexer.New(s, strings.NewReader(s))
	pt := pos.New(pos.WithModel(model))

	pt.Input = l.Output
	pt.Output = output

	go l.Run()
	pt.Run()
}
