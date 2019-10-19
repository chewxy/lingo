package lingo

import (
	"fmt"
	"unicode"
)

//go:generate stringer -type=LexemeType

type LexemeType byte

const (
	EOF LexemeType = iota
	Word
	Disambig
	URI
	Number
	Date
	Time
	Punctuation
	Symbol
	Space
	SystemUse
)

type Lexeme struct {
	Value      string
	LexemeType LexemeType

	Line int
	Col  int
	Pos  int
}

func MakeLexeme(s string, t LexemeType) Lexeme {
	return Lexeme{
		Value:      s,
		LexemeType: t,
		Line:       -1,
		Col:        -1,
		Pos:        -1,
	}
}

func (l Lexeme) Fix() Lexeme {
	if StringIs(l.Value, unicode.IsDigit) {
		l.LexemeType = Number
		return l
	}
	return l
}

func (l Lexeme) String() string {
	switch l.LexemeType {
	case EOF:
		return "EOF"
	default:
		return fmt.Sprintf("%q/%v", l.Value, l.LexemeType)
	}
}

func (l Lexeme) GoString() string {
	switch l.LexemeType {
	case EOF:
		return fmt.Sprintf("EOF: %q (%d, %d, %d)", l.Value, l.Line, l.Col, l.Pos)
	default:
		return fmt.Sprintf("%s: %q (%d, %d, %d)", l.LexemeType, l.Value, l.Line, l.Col, l.Pos)
	}
}

var startLexeme = MakeLexeme("START_LEXEME", SystemUse)
var rootLexeme = MakeLexeme("-ROOT-", SystemUse)
var nullLexeme = MakeLexeme("", SystemUse)

func StartLexeme() Lexeme { return startLexeme }
func RootLexeme() Lexeme  { return rootLexeme }
func NullLexeme() Lexeme  { return nullLexeme }
