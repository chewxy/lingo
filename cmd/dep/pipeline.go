package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/chewxy/lingo"
	"github.com/chewxy/lingo/dep"
	"github.com/chewxy/lingo/lexer"
	"github.com/chewxy/lingo/pos"
)

func receive(deps chan *lingo.Dependency, errs, errChan chan error) {
	defer close(errChan)
	for {
		select {
		case dep, ok := <-deps:
			if !ok {
				continue
			}
			switch *format {
			case "json":
				bs, _ := json.MarshalIndent(dep, "", "\t")
				fmt.Printf("%s\n", string(bs))
			case "dot":
				fmt.Printf("%v\n", dep.Tree().Dot())
			}

		case err := <-errs:
			errChan <- err
		}
	}
}

func pipeline(s string) error {
	lx := lexer.New(s, strings.NewReader(s))
	pt := pos.New(pos.WithModel(POSModel))
	dp := dep.New(DepModel)

	pt.Input = lx.Output
	dp.Input = pt.Output

	errChan := make(chan error)
	go lx.Run()
	go pt.Run()
	go receive(dp.Output, dp.Error, errChan)
	dp.Run()

	return <-errChan
}
