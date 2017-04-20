package treebank

import (
	"strings"
	"testing"

	"github.com/chewxy/lingo"
	"github.com/stretchr/testify/assert"
)

const sampleConllu = `1	President	President	PROPN	NNP	Number=Sing	2	compound	_	_
2	Bush	Bush	PROPN	NNP	Number=Sing	5	nsubj	_	_
3	on	on	ADP	IN	_	4	case	_	_
4	Tuesday	Tuesday	PROPN	NNP	Number=Sing	5	nmod	_	_
5	nominated	nominate	VERB	VBD	Mood=Ind|Tense=Past|VerbForm=Fin	0	root	_	_
6	two	two	NUM	CD	NumType=Card	7	nummod	_	_
7	individuals	individual	NOUN	NNS	Number=Plur	5	dobj	_	_
8	to	to	PART	TO	_	9	mark	_	_
9	replace	replace	VERB	VB	VerbForm=Inf	5	advcl	_	_
10	retiring	retire	VERB	VBG	VerbForm=Ger	11	amod	_	_
11	jurists	jurist	NOUN	NNS	Number=Plur	9	dobj	_	_
12	on	on	ADP	IN	_	14	case	_	_
13	federal	federal	ADJ	JJ	Degree=Pos	14	amod	_	_
14	courts	court	NOUN	NNS	Number=Plur	11	nmod	_	_
15	in	in	ADP	IN	_	18	case	_	_
16	the	the	DET	DT	Definite=Def|PronType=Art	18	det	_	_
17	Washington	Washington	PROPN	NNP	Number=Sing	18	compound	_	_
18	area	area	NOUN	NN	Number=Sing	14	nmod	_	_
19	.	.	PUNCT	.	_	5	punct	_	_

`

func Test_ReadConllu(t *testing.T) {
	assert := assert.New(t)
	st := ReadConllu(strings.NewReader(sampleConllu))[0]

	correctHeads := []int{2, 5, 4, 5, 0, 7, 5, 9, 5, 11, 9, 14, 14, 11, 18, 18, 18, 14, 5}
	assert.Equal(correctHeads, st.Heads)

	// we compare by string to avoid having to build two different test files
	var correctPOS []string
	if lingo.BUILD_TAGSET == "stanfordtags" {
		correctPOS = []string{
			"NNP",
			"NNP",
			"IN",
			"NNP",
			"VBD",
			"CD",
			"NNS",
			"TO",
			"VB",
			"VBG",
			"NNS",
			"IN",
			"JJ",
			"NNS",
			"IN",
			"DT",
			"NNP",
			"NN",
			"FULLSTOP",
		}
	} else {
		correctPOS = []string{
			"PROPN",
			"PROPN",
			"ADP",
			"PROPN",
			"VERB",
			"NUM",
			"NOUN",
			"PART",
			"VERB",
			"VERB",
			"NOUN",
			"ADP",
			"ADJ",
			"NOUN",
			"ADP",
			"DET",
			"PROPN",
			"NOUN",
			"PUNCT",
		}
	}

	assert.Equal(correctPOS, ttos(st.Tags))

	// the stanford tags are not listed in the CONLLU format
	if lingo.BUILD_RELSET != "stanfordrel" {
		var correctRel []string
		correctRel = []string{
			"Compound",
			"NSubj",
			"Case",
			"NMod",
			"Root",
			"NumMod",
			"DObj",
			"Mark",
			"AdvCl",
			"AMod",
			"DObj",
			"Case",
			"AMod",
			"NMod",
			"Case",
			"Det",
			"Compound",
			"NMod",
			"Punct",
		}

		assert.Equal(correctRel, ltos(st.Labels))
	}
}

func ttos(ts []lingo.POSTag) []string {
	retVal := make([]string, len(ts))
	for i, t := range ts {
		retVal[i] = t.String()
	}
	return retVal
}

func ltos(ls []lingo.DependencyType) []string {
	retVal := make([]string, len(ls))
	for i, l := range ls {
		retVal[i] = l.String()
	}
	return retVal
}
