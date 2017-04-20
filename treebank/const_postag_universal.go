// +build !stanfordtags

package treebank

import "github.com/chewxy/lingo"

var posTagTable map[string]lingo.POSTag = map[string]lingo.POSTag{
	"X":     lingo.X,
	"ADJ":   lingo.ADJ,
	"ADP":   lingo.ADP,
	"ADV":   lingo.ADV,
	"AUX":   lingo.AUX,
	"CONJ":  lingo.CONJ,
	"DET":   lingo.DET,
	"INTJ":  lingo.INTJ,
	"NOUN":  lingo.NOUN,
	"NUM":   lingo.NUM,
	"PART":  lingo.PART,
	"PRON":  lingo.PRON,
	"PROPN": lingo.PROPN,
	"PUNCT": lingo.PUNCT,
	"SCONJ": lingo.SCONJ,
	"SYM":   lingo.SYM,
	"VERB":  lingo.VERB,

	"-NULL-":    lingo.X,
	"-ROOT-":    lingo.ROOT_TAG,
	"-UNKNOWN-": lingo.UNKNOWN_TAG,
}
