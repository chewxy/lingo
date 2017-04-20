# lingo #

<img src="https://raw.githubusercontent.com/chewxy/lingo/master/media/gopher_small.png" align="right" />

package `lingo` provides the data structures and algorithms required for natural language processing.

Specifically, it provides a POS Tagger (`lingo/pos`), a Dependency Parser (`lingo/dep`), and a basic tokenizer (`lingo/lexer`) for English. It also provides data structures for holding corpuses (`lingo/corpus`), and treebanks (`lingo/treebank`).

The aim of this package is to provide a production quality pipeline for natural language processing.

# Install #

The package is go-gettable: `go get -u github.com/chewxy/lingo`

This package and its subpackages depend on very few external packages. Here they are:

| Package | Used For | Vitality | Notes | Licence |
|---------|----------|----------|-------|---------|
| [gorgonia](https://github.com/chewxy/gorgonia) | Machine learning | Vital. It won't be hard to rewrite them, but why? | Same author | [Gorgonia Licence](https://github.com/chewxy/gorgonia/blob/master/LICENSE) (Apache 2.0-like) |
| [gographviz](https://github.com/awalterschulze/gographviz) | Visualization of annotations, and other graph-related visualizations | Vital for visualizations, which are a nice-to-have feature | API last changed 12th April 2017 | [gographviz licence](https://github.com/awalterschulze/gographviz/blob/master/LICENSE) (Apache 2.0) |
| [errors](https://github.com/pkg/errors)  | Errors   | The package won't die without it, but it's a very nice to have | Stable API for the past year | [errors licence](https://github.com/pkg/errors/blob/master/LICENSE) (MIT/BSD like) |
| [set](https://github.com/xtgo/set) | Set operations | Can be easily replaced | Stable API for the past year | [set licence](https://github.com/xtgo/set/blob/master/LICENSE) (MIT/BSD-like) |

# Usage #

See the individual packages for usage. There is also a bunch of executables in the `cmd` directory. They're meant to be examples as to how a natural language processing pipeline can be set up.

A natural language pipeline with this package is heavily channels driven. Here's is an example for dependency parsing:

```go
func main() {
	inputString: `The cat sat on the mat`
	lx := lexer.New("dummy", strings.NewReader(inputString)) // lexer - required to break a sentence up into words. 
	pt := pos.New(pos.WithModel(posModel))                   // POS Tagger - required to tag the words with a part of speech tag.
	dp := dep.New(depModel)                                  // Creates a new parser

	// set up a pipeline
	pt.Input = lx.Output
	dp.Input = pt.Output

	// run all
	go lx.Run()
	go pt.Run()
	go dp.Run()

	// wait to receive:
	for {
		select {
		case d := <- dp.Output:
			// do something
		case err:= <-dp.Error:
			// handle error
		}
	}

}

```



# How It Works #
For specific tasks (POS tagging, parsing, named entity recognition etc), refer to the README of each subpackage. This package on its own mainly provides the data structures that the subpackages will use.

Perhaps the most important data structure is the `*Annotation` structure. It basically holds a word and the associated metadata for the word. 

For dependency parses, the graph takes three forms: `*Dependency`, `*DependencyTree` and `*Annotation`. All three forms are convertable from one to another. TODO: explain rationale behind each data type.

## Quirks ##

### Very Oddly Specific POS Tags and Dependency Rel Types ###

A particular quirk you may have noticed is that the `POSTag` and `DependencyType` are hard coded in as constants. This package does in fact provide two variations of each: one from Stanford/Penn Treebank and one from [UniversalDependencies](http://universaldependencies.org/). 

The main reason for hardcoding these are mainly for performance reasons - knowing ahead how much to allocate reduces a lot of additional work the program has to do. It also reduces the chances of mutating a global variable. 

Of course this comes as a tradeoff - programs are limited to these two options. Thankfully there are only a limited number of POS Tag and Dependency Relation types. Two of the most popular ones (Stanford/PTB and Universal Dependencies) have been implemented. 

The following build tags are supported:

* stanfordtags
* universaltags 
* stanfordrel
* universalrel

To use a specific tagset or relset, build your program thusly: `go build -tags='stanfordtags'`. 

The default tag and dependency rel types are the universal dependencies version.

### Lexer ###

You should also note that the tokenizer, `lingo/lexer` is not your usual run-of-the-mill NLP tokenizer. It's a tokenizer that tokenizes by space, with some specific rules for English. It was inspired by Rob Pike's talk on lexers. I thought it'd be cool to write something like that for NLP.

The test cases in package `lingo/lexer` showcases how it handles unicode, and other pathalogical english.

# Contributing #
see CONTRIBUTING.md for more info

# Licence #

This package is licenced under the MIT licence.
