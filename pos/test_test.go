package pos

import (
	"github.com/chewxy/lingo"
	"github.com/kljensen/snowball"
)

type dummyLem struct{}

func (dummyLem) Lemmatize(s string, pt lingo.POSTag) ([]string, error) {
	if len(s) > 3 {
		return []string{
			s[:2],
		}, nil
	}
	return []string{""}, nil
}

type dummyStemmer struct{}

func (dummyStemmer) Stem(s string) (string, error) {
	return snowball.Stem(s, "english", true)
}

var clusters = map[string]lingo.Cluster{
	"TEst": 1,
	"Test": 1,
	"test": 1,
}

type dummyFix struct {
	dummyStemmer
	dummyLem
}

func (dummyFix) Clusters() (map[string]lingo.Cluster, error) { return clusters, nil }

const conllu = `1	From	from	ADP	IN	_	3	case	_	_
2	the	the	DET	DT	Definite=Def|PronType=Art	3	det	_	_
3	AP	AP	PROPN	NNP	Number=Sing	4	nmod	_	_
4	comes	come	VERB	VBZ	Mood=Ind|Number=Sing|Person=3|Tense=Pres|VerbForm=Fin	0	root	_	_
5	this	this	DET	DT	Number=Sing|PronType=Dem	6	det	_	_
6	story	story	NOUN	NN	Number=Sing	4	nsubj	_	_
7	:	:	PUNCT	:	_	4	punct	_	_

1	President	President	PROPN	NNP	Number=Sing	2	compound	_	_
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

1	Bush	Bush	PROPN	NNP	Number=Sing	2	nsubj	_	_
2	nominated	nominate	VERB	VBD	Mood=Ind|Tense=Past|VerbForm=Fin	0	root	_	_
3	Jennifer	Jennifer	PROPN	NNP	Number=Sing	5	compound	_	_
4	M.	M.	PROPN	NNP	Number=Sing	5	compound	_	_
5	Anderson	Anderson	PROPN	NNP	Number=Sing	2	dobj	_	_
6	for	for	ADP	IN	_	11	case	_	_
7	a	a	DET	DT	Definite=Ind|PronType=Art	11	det	_	_
8	15	15	NUM	CD	NumType=Card	10	nummod	_	_
9	-	-	PUNCT	HYPH	_	10	punct	_	_
10	year	year	NOUN	NN	Number=Sing	11	compound	_	_
11	term	term	NOUN	NN	Number=Sing	2	nmod	_	_
12	as	as	ADP	IN	_	14	case	_	_
13	associate	associate	ADJ	JJ	Degree=Pos	14	amod	_	_
14	judge	judge	NOUN	NN	Number=Sing	11	nmod	_	_
15	of	of	ADP	IN	_	18	case	_	_
16	the	the	DET	DT	Definite=Def|PronType=Art	18	det	_	_
17	Superior	Superior	PROPN	NNP	Number=Sing	18	compound	_	_
18	Court	Court	PROPN	NNP	Number=Sing	14	nmod	_	_
19	of	of	ADP	IN	_	21	case	_	_
20	the	the	DET	DT	Definite=Def|PronType=Art	21	det	_	_
21	District	District	PROPN	NNP	Number=Sing	18	nmod	_	_
22	of	of	ADP	IN	_	23	case	_	_
23	Columbia	Columbia	PROPN	NNP	Number=Sing	21	nmod	_	_
24	,	,	PUNCT	,	_	2	punct	_	_
25	replacing	replace	VERB	VBG	VerbForm=Ger	2	advcl	_	_
26	Steffen	Steffen	PROPN	NNP	Number=Sing	28	compound	_	_
27	W.	W.	PROPN	NNP	Number=Sing	28	compound	_	_
28	Graae	Graae	PROPN	NNP	Number=Sing	25	dobj	_	_
29	.	.	PUNCT	.	_	2	punct	_	_

1	We	we	PRON	PRP	Case=Nom|Number=Plur|Person=1|PronType=Prs	3	nsubj	_	_
2	've	have	AUX	VBP	Mood=Ind|Tense=Pres|VerbForm=Fin	3	aux	_	_
3	grown	grow	VERB	VBN	Tense=Past|VerbForm=Part	0	root	_	_
4	up	up	ADP	RP	_	3	compound:prt	_	_
5	.	.	PUNCT	.	_	3	punct	_	_`
