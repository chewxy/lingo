package treebank

import "github.com/chewxy/lingo"

var alreadyLogged map[string]bool = make(map[string]bool)

// TODO : CHECK
func StringToLexType(tag string) lingo.LexemeType {
	var lexType lingo.LexemeType
	switch tag {
	case "NUM":
		lexType = lingo.Number
	case "PUNCT":
		lexType = lingo.Punctuation
	case "SYM":
		lexType = lingo.Symbol
	default:
		lexType = lingo.Word
	}
	return lexType
}

func StringToPOSTag(tag string) (lingo.POSTag, bool) {
	t, ok := posTagTable[tag]

	return t, ok
}

func StringToDependencyType(ud string) (lingo.DependencyType, bool) {
	dt, ok := dependencyTable[ud]

	return dt, ok
}

func reset() (lingo.LexemeSentence, []lingo.POSTag, []int, []lingo.DependencyType) {
	s := lingo.NewLexemeSentence()
	st := make([]lingo.POSTag, 0)
	sh := make([]int, 0)
	sdt := make([]lingo.DependencyType, 0)

	return s, st, sh, sdt
}

func finish(s lingo.LexemeSentence, st []lingo.POSTag, sh []int, sdt []lingo.DependencyType, sentences []SentenceTag) []SentenceTag {
	sentenceTag := SentenceTag{s, st, sh, sdt}
	sentences = append(sentences, sentenceTag)

	return sentences
}
