package dep

import (
	"bufio"
	"crypto/md5"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/chewxy/lingo"
	"github.com/chewxy/lingo/treebank"
	"github.com/kljensen/snowball"
)

type dummyLem struct{}

func (dummyLem) Lemmatize(s string, pt lingo.POSTag) ([]string, error) {
	return nil, componentUnavailable("lemmatizer")
}

type dummyStemmer struct{}

func (dummyStemmer) Stem(s string) (string, error) {
	return snowball.Stem(s, "english", true)
}

type dummyFix struct {
	dummyStemmer
	dummyLem
}

func (dummyFix) Clusters() (map[string]lingo.Cluster, error) {
	return nil, componentUnavailable("clusters")
}

const nnps = `1	Guerrillas	guerrilla	NOUN	NNS	Number=Plur	2	nsubj	_	_
2	threatened	threaten	VERB	VBD	Mood=Ind|Tense=Past|VerbForm=Fin	0	root	_	_
3	to	to	PART	TO	_	4	mark	_	_
4	assassinate	assassinate	VERB	VB	VerbForm=Inf	2	xcomp	_	_
5	Prime	Prime	PROPN	NNP	Number=Sing	6	compound	_	_
6	Minister	Minister	PROPN	NNP	Number=Sing	8	compound	_	_
7	Iyad	Iyad	PROPN	NNP	Number=Sing	8	compound	_	_
8	Allawi	Allawi	PROPN	NNP	Number=Sing	4	dobj	_	_
9	and	and	CONJ	CC	_	8	cc	_	_
10	Minister	Minister	PROPN	NNP	Number=Sing	14	compound	_	_
11	of	of	ADP	IN	_	12	case	_	_
12	Defense	Defense	PROPN	NNP	Number=Sing	10	nmod	_	_
13	Hazem	Hazem	PROPN	NNP	Number=Sing	14	compound	_	_
14	Shaalan	Shaalan	PROPN	NNP	Number=Sing	8	conj	_	_
15	in	in	ADP	IN	_	16	case	_	_
16	retaliation	retaliation	NOUN	NN	Number=Sing	4	nmod	_	_
17	for	for	ADP	IN	_	19	case	_	_
18	the	the	DET	DT	Definite=Def|PronType=Art	19	det	_	_
19	attack	attack	NOUN	NN	Number=Sing	16	nmod	_	_
20	.	.	PUNCT	.	_	2	punct	_	_

`
const simple = `1	Yet	yet	CONJ	CC	_	5	cc	_	_
2	we	we	PRON	PRP	Case=Nom|Number=Plur|Person=1|PronType=Prs	5	nsubj	_	_
3	did	do	AUX	VBD	Mood=Ind|Tense=Past|VerbForm=Fin	5	aux	_	_
4	n't	not	PART	RB	_	5	neg	_	_
5	charge	charge	VERB	VB	VerbForm=Inf	0	root	_	_
6	them	they	PRON	PRP	Case=Acc|Number=Plur|Person=3|PronType=Prs	5	dobj	_	_
7	for	for	ADP	IN	_	9	case	_	_
8	the	the	DET	DT	Definite=Def|PronType=Art	9	det	_	_
9	evacuation	evacuation	NOUN	NN	Number=Sing	5	nmod	_	_
10	.	.	PUNCT	.	_	5	punct	_	_

`

const med = `1	President	President	PROPN	NNP	Number=Sing	2	compound	_	_
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

const long = `1	Now	now	ADV	RB	_	5	advmod	_	_
2	,	,	PUNCT	,	_	5	punct	_	_
3	I	I	PRON	PRP	Case=Nom|Number=Sing|Person=1|PronType=Prs	5	nsubj	_	_
4	would	would	AUX	MD	VerbForm=Fin	5	aux	_	_
5	argue	argue	VERB	VB	VerbForm=Inf	0	root	_	_
6	that	that	SCONJ	IN	_	11	mark	_	_
7	one	one	PRON	PRP	_	11	nsubj	_	_
8	could	could	AUX	MD	VerbForm=Fin	11	aux	_	_
9	have	have	AUX	VB	VerbForm=Inf	11	aux	_	_
10	reasonably	reasonably	ADV	RB	_	11	advmod	_	_
11	predicted	predict	VERB	VBN	Tense=Past|VerbForm=Part	5	ccomp	_	_
12	that	that	SCONJ	IN	_	19	mark	_	_
13	some	some	DET	DT	_	14	det	_	_
14	form	form	NOUN	NN	Number=Sing	19	nsubj	_	_
15	of	of	ADP	IN	_	17	case	_	_
16	military	military	ADJ	JJ	Degree=Pos	17	amod	_	_
17	violence	violence	NOUN	NN	Number=Sing	14	nmod	_	_
18	was	be	VERB	VBD	Mood=Ind|Number=Sing|Person=3|Tense=Past|VerbForm=Fin	19	cop	_	_
19	likely	likely	ADJ	JJ	Degree=Pos	11	ccomp	_	_
20	to	to	PART	TO	_	21	mark	_	_
21	occur	occur	VERB	VB	VerbForm=Inf	19	xcomp	_	_
22	in	in	ADP	IN	_	23	case	_	_
23	Lebanon	Lebanon	PROPN	NNP	Number=Sing	21	nmod	_	_
24	-LRB-	-lrb-	PUNCT	-LRB-	_	25	punct	_	_
25	considering	consider	VERB	VBG	VerbForm=Ger	19	advcl	_	_
26	that	that	SCONJ	IN	_	31	mark	_	_
27	the	the	DET	DT	Definite=Def|PronType=Art	28	det	_	_
28	country	country	NOUN	NN	Number=Sing	31	nsubj	_	_
29	has	have	AUX	VBZ	Mood=Ind|Number=Sing|Person=3|Tense=Pres|VerbForm=Fin	31	aux	_	_
30	been	be	AUX	VBN	Tense=Past|VerbForm=Part	31	aux	_	_
31	experiencing	experience	VERB	VBG	Tense=Pres|VerbForm=Part	25	ccomp	_	_
32	some	some	DET	DT	_	33	det	_	_
33	form	form	NOUN	NN	Number=Sing	31	dobj	_	_
34	of	of	ADP	IN	_	35	case	_	_
35	conflict	conflict	NOUN	NN	Number=Sing	33	nmod	_	_
36	for	for	ADP	IN	_	41	case	_	_
37	approximately	approximately	ADV	RB	_	41	advmod	_	_
38	the	the	DET	DT	Definite=Def|PronType=Art	41	det	_	_
39	last	last	ADJ	JJ	Degree=Pos	41	amod	_	_
40	32	32	NUM	CD	NumType=Card	41	nummod	_	_
41	years	year	NOUN	NNS	Number=Plur	31	nmod	_	_
42	-RRB-	-rrb-	PUNCT	-RRB-	_	25	punct	_	_
43	.	.	PUNCT	.	_	5	punct	_	_

`

const cvconllu = `1	Google	Google	PROPN	NNP	Number=Sing	6	nsubj	_	_
2	is	be	VERB	VBZ	Mood=Ind|Number=Sing|Person=3|Tense=Pres|VerbForm=Fin	6	cop	_	_
3	a	a	DET	DT	Definite=Ind|PronType=Art	6	det	_	_
4	nice	nice	ADJ	JJ	Degree=Pos	6	amod	_	_
5	search	search	NOUN	NN	Number=Sing	6	compound	_	_
6	engine	engine	NOUN	NN	Number=Sing	0	root	_	_
7	.	.	PUNCT	.	_	6	punct	_	_

1	Does	do	AUX	VBZ	Mood=Ind|Number=Sing|Person=3|Tense=Pres|VerbForm=Fin	3	aux	_	_
2	anybody	anybody	NOUN	NN	Number=Sing	3	nsubj	_	_
3	use	use	VERB	VB	VerbForm=Inf	0	root	_	_
4	it	it	PRON	PRP	Case=Acc|Gender=Neut|Number=Sing|Person=3|PronType=Prs	3	dobj	_	_
5	for	for	ADP	IN	_	6	case	_	_
6	anything	anything	NOUN	NN	Number=Sing	3	nmod	_	_
7	else	else	ADJ	JJ	Degree=Pos	6	amod	_	_
8	?	?	PUNCT	.	_	3	punct	_	_

`

func lotsaNNP() *lingo.Dependency {
	readr := strings.NewReader(nnps)
	sentenceTags := treebank.ReadConllu(readr)

	return sentenceTags[0].Dependency(dummyFix{})
}

// simpleSentence has 10 words
func simpleSentence() []treebank.SentenceTag {
	readr := strings.NewReader(simple)
	return treebank.ReadConllu(readr)
}

func mediumSentence() []treebank.SentenceTag {
	readr := strings.NewReader(med)
	return treebank.ReadConllu(readr)
}

// longSentence has 44 words
func longSentence() []treebank.SentenceTag {
	readr := strings.NewReader(long)
	return treebank.ReadConllu(readr)
}

func allSentences() []treebank.SentenceTag {
	sentenceTags := treebank.ReadConllu(strings.NewReader(nnps))
	sentenceTags = append(sentenceTags, treebank.ReadConllu(strings.NewReader(simple))...)
	sentenceTags = append(sentenceTags, treebank.ReadConllu(strings.NewReader(med))...)
	sentenceTags = append(sentenceTags, treebank.ReadConllu(strings.NewReader(long))...)
	return sentenceTags
}

func cvSentences() []treebank.SentenceTag {
	return treebank.ReadConllu(strings.NewReader(cvconllu))
}

func hash(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func cache(input string, s lingo.AnnotatedSentence) {
	hashfilename := "cached/" + hash(input) + ".cached"
	f, err := os.Create(hashfilename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	defer w.Flush()

	encoder := gob.NewEncoder(w)

	if err := encoder.Encode(s); err != nil {
		log.Fatal(err)
	}
}

func useCached(filename string) *lingo.Dependency {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	r := bufio.NewReader(f)
	decoder := gob.NewDecoder(r)

	var sentence lingo.AnnotatedSentence
	if err := decoder.Decode(&sentence); err != nil {
		log.Fatal(err)
	}
	// fixes ID and what nots
	sentence.Fix()

	dep := sentence.Dependency()
	return dep
}
