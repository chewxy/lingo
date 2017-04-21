package main

import (
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/abiosoft/ishell"
	"github.com/chewxy/lingo"
	"github.com/pkg/browser"
)

func main() {
	io()
	shell := ishell.New()

	var d *lingo.Dependency
	// var sent lingo.AnnotatedSentence
	var err error
	shell.AddCmd(&ishell.Cmd{
		Name: "dep",
		Help: "perform dependency parsing",
		Func: func(c *ishell.Context) {
			c.ShowPrompt(false)
			defer c.ShowPrompt(true)

			c.Print("Query: ")
			query := c.ReadLine()

			if d, err = pipeline(query); err != nil {
				c.Printf("Error: %v", err)
			}

			c.Printf("%v\n", d)
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "show",
		Help: "show dependency parse on browser",
		Func: func(c *ishell.Context) {
			var tmp *os.File
			if tmp, err = ioutil.TempFile("", "dep"); err != nil {
				c.Printf("Cannot open file %v\n", err)
				return
			}
			defer os.Remove(tmp.Name())

			c.Printf("%v\n", tmp.Name())

			dot := d.Tree().Dot()
			tmp.Write([]byte(dot))
			if err := tmp.Close(); err != nil {
				c.Printf("Error closing file %v", err)
			}
			cmd := exec.Command("dot", "-Tpng", "-O", tmp.Name())
			if err = cmd.Run(); err != nil {
				c.Printf("Cannot execute dot: %v\n", err)
			}

			browser.OpenFile(tmp.Name() + ".png")

		},
	})
	shell.Start()
}
