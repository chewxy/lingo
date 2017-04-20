package corpus

import (
	"strings"

	"github.com/chewxy/lingo/treebank"
)

const sample1Gram = `the	23135851162
of	13151942776
and	12997637966
to	12136980858
a	9081174698
in	8469404971
for	5933321709`

func mediumSentence() []treebank.SentenceTag {
	conllu := `1	President	President	PROPN	NNP	Number=Sing	2	compound	_	_
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

	readr := strings.NewReader(conllu)
	return treebank.ReadConllu(readr)
}

const EPSILON64 float64 = 1e-10

func floatEquals64(a, b float64) bool {
	if (a-b) < EPSILON64 && (b-a) < EPSILON64 {
		return true
	}
	return false
}
