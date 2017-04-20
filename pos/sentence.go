package pos

import "github.com/chewxy/lingo"

// "log"

func (p *Tagger) getSentences() {
	defer close(p.sentences)

	var sentence lingo.AnnotatedSentence
	sentence = append(sentence, lingo.RootAnnotation())

	for lexeme := range p.Input {
		if lexeme.LexemeType != lingo.EOF {
			a := lingo.NewAnnotation()
			a.Lexeme = lexeme
			if err := a.Process(p); err != nil {
				panic(err) // for now
			}
			sentence = append(sentence, a)
		} else {
			p.sentences <- sentence

			// reset
			sentence = lingo.AnnotatedSentence{lingo.RootAnnotation()}
		}

		// TODO: Sentence splitting
	}
}
