package lingo

import (
	"fmt"
	"strings"
	"unicode"
)

// WordFlags represent the types a word may be. A word may have multiple flags
type WordFlag uint32

const (
	NoFlag WordFlag = iota
	IsLetter
	IsAscii
	IsDigit
	IsLower
	IsPunct
	IsSpace
	IsTitle
	IsUpper
	LikeURL
	LikeNum
	LikeEmail
	IsStopWord
	IsOOV // for ner

	MAXFLAG
)

func (f WordFlag) String() string {
	return fmt.Sprintf("%014b", f)
}

func (l Lexeme) Flags() WordFlag {
	var wf WordFlag

	s := l.Value

	if StringIs(s, unicode.IsLetter) {
		wf |= (1 << IsLetter)
	}

	if StringIs(s, unicode.IsDigit) {
		wf |= (1 << IsDigit)
	}

	if StringIs(s, isAscii) {
		wf |= (1 << IsAscii)
	}

	if StringIs(s, unicode.IsLower) {
		wf |= (1 << IsLower)
	}

	if StringIs(s, unicode.IsPunct) {
		wf |= (1 << IsPunct)
	}

	if StringIs(s, unicode.IsSpace) {
		wf |= (1 << IsSpace)
	}

	if StringIs(s, unicode.IsUpper) {
		wf |= (1 << IsUpper)
	}

	if l.LexemeType == URI {
		wf |= (1 << LikeURL)
	}

	if _, ok := NumberWords[strings.ToLower(s)]; ok {
		wf |= (1 << LikeNum)
	}

	if _, ok := stopwords[s]; ok {
		wf |= (1 << IsStopWord)
	}

	if len(s) > 0 {
		if (unicode.IsUpper(rune(s[0])) || unicode.IsTitle(rune(s[0]))) && StringIs(s[1:], unicode.IsLower) {
			wf |= (1 << IsTitle)
		}
	}

	return wf
}
