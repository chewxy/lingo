package treebank

import (
	"math/rand"

	"github.com/chewxy/lingo"
)

// SentenceTag is a struc that holds a sentence, tags, heads and labels
type SentenceTag struct {
	Sentence lingo.LexemeSentence
	Tags     []lingo.POSTag
	Heads    []int
	Labels   []lingo.DependencyType
}

func (s SentenceTag) AnnotatedSentence(f lingo.AnnotationFixer) lingo.AnnotatedSentence {
	retVal := lingo.NewAnnotatedSentence()
	retVal = append(retVal, lingo.RootAnnotation())

	for i, lex := range s.Sentence {
		a := lingo.NewAnnotation()
		a.Lexeme = lex
		a.POSTag = s.Tags[i]
		a.DependencyType = s.Labels[i]

		// should panic, because SentenceTag is only ever used during training
		if err := a.Process(f); err != nil {
			panic(err)
		}

		retVal = append(retVal, a)
	}

	// add heads
	for i, a := range retVal {
		if i == 0 {
			continue
		}
		a.SetHead(retVal[s.Heads[i-1]])
	}

	retVal.Fix()

	return retVal
}

func (s SentenceTag) Dependency(f lingo.AnnotationFixer) *lingo.Dependency {
	sentence := s.AnnotatedSentence(f)
	dep := sentence.Dependency()

	return dep
}

func (s SentenceTag) String() string {
	return s.Sentence.String()
}

func ShuffleSentenceTag(s []SentenceTag) []SentenceTag {
	rand.Seed(1337)
	for i := range s {
		j := rand.Intn(i + 1)
		s[i], s[j] = s[j], s[i]
	}

	return s
}

/* UTILITY FUNCTIONS */

func WrapLexemeSentence(sentence lingo.LexemeSentence) lingo.LexemeSentence {
	retSentence := lingo.NewLexemeSentence()
	retSentence = append(retSentence, lingo.StartLexeme())
	retSentence = append(retSentence, sentence...)
	retSentence = append(retSentence, lingo.RootLexeme())
	return retSentence
}

func WrapTags(tagList []lingo.POSTag) []lingo.POSTag {
	retVal := append([]lingo.POSTag{lingo.X}, tagList...)
	retVal = append(retVal, lingo.X)
	return retVal
}

func WrapHeads(heads []int) []int {
	retVal := append([]int{0}, heads...)
	retVal = append(retVal, 0)
	return retVal
}

func WrapDeps(deps []lingo.DependencyType) []lingo.DependencyType {
	retVal := append([]lingo.DependencyType{lingo.Dep}, deps...)
	retVal = append(retVal, lingo.Dep)
	return retVal
}
