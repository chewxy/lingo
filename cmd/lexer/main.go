package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/chewxy/lingo"
	"github.com/chewxy/lingo/lexer"
)

var input = flag.String("input", "", "input string to lex")
var output = make(chan lingo.Lexeme)

func receieve() {
	for l := range output {
		fmt.Printf("%v\n", l)
	}
}

func main() {
	flag.Parse()

	s := *input

	go receieve()
	l := lexer.New(s, strings.NewReader(s))
	l.Output = output
	l.Run()
}
