package pos

import (
	"log"
	"strings"
	"testing"

	"github.com/chewxy/lingo"
	"github.com/chewxy/lingo/lexer"
	"github.com/chewxy/lingo/treebank"
)

func TestEverything(t *testing.T) {
	sentences := treebank.ReadConllu(strings.NewReader(conllu))

	sentence := "President Bush comes on federal courts."

	p := New(WithCluster(clusters), WithLemmatizer(dummyLem{}), WithStemmer(dummyStemmer{}))
	p.Train(sentences, 200)

	l := lexer.New(sentence, strings.NewReader(sentence))
	p2 := p.Clone()
	p2.Input = l.Output

	var correct string
	if lingo.BUILD_TAGSET == "stanfordtags" {
		correct = "-ROOT-/ROOT_TAG President/NNP Bush/NNP comes/DT on/IN federal/JJ courts/NN ./FULLSTOP"
	} else {
		correct = "-ROOT-/ROOT_TAG President/PROPN Bush/PROPN comes/VERB on/ADP federal/ADJ courts/NOUN ./PUNCT"
	}

	go l.Run()
	go p2.Run()
	for a := range p2.Output {

		// this clearly isn't gonna be accurate, given the stubbed out Lemmatizer
		if a.String() != correct {
			t.Error("Something went wrong with the POSTagging")
			log.Printf("%v", a)
		}
	}

}
