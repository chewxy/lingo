package lexer

import (
	"strings"
	"testing"

	"github.com/chewxy/lingo"
)

type lexerTest struct {
	name string
	s    string

	lexemes []lingo.Lexeme
}

var lexerTests = []lexerTest{
	// {"empty", "", []lingo.Lexeme{
	// 	{"", lingo.EOF, 0, 1, 0},
	// }},
	//
	// {"spaces", " \t", []lingo.Lexeme{
	// 	{"", lingo.EOF, 0, 3, 2},
	// }},
	//
	// {"newlines", "\n\r\n\n", []lingo.Lexeme{
	// 	{"", lingo.EOF, 3, 5, 4},
	// }},
	//
	// {"simple text", "hello world", []lingo.Lexeme{
	// 	{"hello", lingo.Word, 0, 1, 0},
	// 	{"world", lingo.Word, 0, 7, 6},
	// 	{"", lingo.EOF, 0, 12, 11},
	// }},
	//
	// {"simple number", "3.1415", []lingo.Lexeme{
	// 	{"3.1415", lingo.Number, 0, 1, 0},
	// 	{"", lingo.EOF, 0, 12, 5},
	// }},

	{"advanced numerology", "3.14 -1.618", []lingo.Lexeme{
		{"3.14", lingo.Number, 0, 1, 0},
		{"-1.618", lingo.Number, 0, 6, 5},
		{"", lingo.EOF, 0, 11, 10},
	}},

	// {"advanced numerology", "3.14 -1.618 6.023e23 1e-13", []lingo.Lexeme{
	// 	{"3.14", lingo.Number, 0, 1, 0},
	// 	{"-1.618", lingo.Number, 0, 6, 5},
	// 	{"6.023e23", lingo.Number, 0, 13, 12},
	// 	{"1e-13", lingo.Number, 0, 21, 20},
	// 	{"", lingo.EOF, 0, 26, 25},
	// }},
	//
	// {"esoteric numerology", "1/2 1 1/4", []lingo.Lexeme{
	// 	{"1/2", lingo.Number, 0, 1, 0},
	// 	{"1", lingo.Number, 0, 5, 4},
	// 	{"1/4", lingo.Number, 0, 7, 6},
	// 	{"", lingo.EOF, 0, 10, 9},
	// }},
	//
	// {"text with numbers", "one plus 1 don't equals 3", []lingo.Lexeme{
	// 	{"one", lingo.Word, 0, 1, 0},
	// 	{"plus", lingo.Word, 0, 5, 4},
	// 	{"1", lingo.Number, 0, 10, 9},
	// 	{"do", lingo.Word, 0, 12, 11},
	// 	{"n't", lingo.Word, 0, 14, 13},
	// 	{"equals", lingo.Word, 0, 18, 17},
	// 	{"3", lingo.Number, 0, 24, 23},
	// 	{"", lingo.EOF, 0, 25, 24},
	// }},
	//
	// {"text with numbers + punct", "First111!.!", []lingo.Lexeme{
	// 	{"First111", lingo.Word, 0, 1, 0},
	// 	{"!.!", lingo.Punctuation, 0, 9, 8},
	// 	{"", lingo.EOF, 0, 10, 9},
	// }},
	//
	// {"text with verb contractions", "You're panic'd I'll get'em I've", []lingo.Lexeme{
	// 	{"You", lingo.Word, 0, 1, 0},
	// 	{"'re", lingo.Word, 0, 3, 2},
	// 	{"panic", lingo.Word, 0, 8, 7},
	// 	{"'d", lingo.Word, 0, 13, 12},
	// 	{"I", lingo.Word, 0, 16, 15},
	// 	{"'ll", lingo.Word, 0, 17, 16},
	// 	{"get", lingo.Word, 0, 21, 20},
	// 	{"'em", lingo.Word, 0, 24, 23},
	// 	{"I", lingo.Word, 0, 27, 26},
	// 	{"'ve", lingo.Word, 0, 30, 29},
	// 	{"", lingo.EOF, 0, 33, 32},
	// }},
	//
	// {"email", "dont@email.this", []lingo.Lexeme{
	// 	{"dont@email.this", lingo.Word, 0, 1},
	// 	{"", lingo.EOF, 0, 10},
	// }},
	//
	// {"plain dashes should not be numbers", "this case - like so", []lingo.Lexeme{
	// 	{"this", lingo.Word, 0, 1},
	// 	{"case", lingo.Word, 0, 5},
	// 	{"-", lingo.Punctuation, 0, 6},
	// 	{"like", lingo.Word, 0, 8},
	// 	{"so", lingo.Word, 0, 13},
	// 	{"", lingo.EOF, 0, 14},
	// }},
	//
	// {"parens should be printed", "like (this)", []lingo.Lexeme{
	// 	{"like", lingo.Word, 0, 1},
	// 	{"(", lingo.Punctuation, 0, 5},
	// 	{"this", lingo.Word, 0, 6},
	// 	{")", lingo.Punctuation, 0, 10},
	// 	{"", lingo.EOF, 0, 11},
	// }},
	//
	// {"parenthesis should be considered separate", "USA(United States of America)", []lingo.Lexeme{
	// 	{"USA", lingo.Word, 0, 1},
	// 	{"(", lingo.Punctuation, 0, 1},
	// 	{"United", lingo.Word, 0, 1},
	// 	{"States", lingo.Word, 0, 1},
	// 	{"of", lingo.Word, 0, 1},
	// 	{"America", lingo.Word, 0, 1},
	// 	{")", lingo.Punctuation, 0, 1},
	// 	{"", lingo.EOF, 0, 0},
	// }},
	//
	// {"midstream puncuation", "like:this", []lingo.Lexeme{
	// 	{"like", lingo.Word, 0, 1},
	// 	{":", lingo.Punctuation, 0, 5},
	// 	{"this", lingo.Word, 0, 6},
	// 	{"", lingo.EOF, 0, 7},
	// }},
	//
	// {"midstream symbols", "e-meet ke$ha by e-mail $ell anti-inflammatory", []lingo.Lexeme{
	// 	{"e-meet", lingo.Word, 0, 1},
	// 	{"ke$ha", lingo.Word, 0, 1},
	// 	{"by", lingo.Word, 0, 1},
	// 	{"e-mail", lingo.Word, 0, 1},
	// 	{"$", lingo.Symbol, 0, 1},
	// 	{"ell", lingo.Word, 0, 1},
	// 	{"anti-inflammatory", lingo.Word, 0, 1},
	// 	{"", lingo.EOF, 0, 0},
	// }},
	//
	// {"abbrev", "USB, made in U.S.A. e.g t/away c/o", []lingo.Lexeme{
	// 	{"USB", lingo.Word, 0, 1},
	// 	{",", lingo.Punctuation, 0, 4},
	// 	{"made", lingo.Word, 0, 6},
	// 	{"in", lingo.Word, 0, 11},
	// 	{"U.S.A", lingo.Word, 0, 14},
	// 	{".", lingo.Punctuation, 0, 19},
	// 	{"e.g", lingo.Word, 0, 0},
	// 	{"t/away", lingo.Word, 0, 0},
	// 	{"c/o", lingo.Word, 0, 0},
	// 	{"", lingo.EOF, 0, 20},
	// }},
	//
	// {"date time", "1970/1/1 00:00 00:00:00", []lingo.Lexeme{
	// 	{"1970/1/1", lingo.Date, 0, 1},
	// 	{"00:00", lingo.Time, 0, 1},
	// 	{"00:00:00", lingo.Time, 0, 20},
	// 	{"", lingo.EOF, 0, 20},
	// }},
	//
	// {"date time with dashes", "31-12-1970", []lingo.Lexeme{
	// 	{"31/12/1970", lingo.Date, 0, 1},
	// 	{"", lingo.EOF, 0, 11},
	// }},
	//
	// {"URI", "wobsite: http://www.wobsite.something.this/is/still/a.url", []lingo.Lexeme{
	// 	{"wobsite", lingo.Word, 0, 1},
	// 	{":", lingo.Punctuation, 0, 8},
	// 	{"http://www.wobsite.something.this/is/still/a.url", lingo.URI, 0, 10},
	// 	{"", lingo.EOF, 0, 20},
	// }},
	//
	// {"proper sentence", "hello world.", []lingo.Lexeme{
	// 	{"hello", lingo.Word, 0, 1},
	// 	{"world", lingo.Word, 0, 6},
	// 	{".", lingo.Punctuation, 0, 7},
	// 	{"", lingo.EOF, 0, 8},
	// }},
	//
	// // Naive and Cafe uses combination diacritics, while the rest are just unicode
	// // The lexer should normalize all the things
	// {"pathological english words", "Façade à la Naïve Château Café", []lingo.Lexeme{
	// 	{"Façade", lingo.Word, 0, 1},
	// 	{"à", lingo.Word, 0, 1},
	// 	{"la", lingo.Word, 0, 1},
	// 	{"Naïve", lingo.Word, 0, 1},
	// 	{"Château", lingo.Word, 0, 1},
	// 	{"Café", lingo.Word, 0, 1},
	// 	{"", lingo.EOF, 0, 0},
	// }},
	//
	// // just plain fucked
	// {"jpf", "你好 العالم", []lingo.Lexeme{
	// 	{"你好", lingo.Word, 0, 1},
	// 	{"العالم", lingo.Word, 0, 1},
	// 	{"", lingo.EOF, 0, 0},
	// }},
}

func testLexer(lts *lexerTest) []lingo.Lexeme {
	l := New(lts.name, strings.NewReader(lts.s))
	var retVal []lingo.Lexeme

	go l.Run()
	for lex := range l.Output {
		retVal = append(retVal, lex)
	}
	return retVal
}

func TestLexer(t *testing.T) {
	for _, lts := range lexerTests {
		lexemes := testLexer(&lts)

		if len(lts.lexemes) != len(lexemes) {
			t.Errorf("Test %q: Expected %d lexemes. Got %d instead: %v", lts.name, len(lts.lexemes), len(lexemes), lexemes)
			continue
		}

		for i, lex := range lexemes {
			if lex.LexemeType != lts.lexemes[i].LexemeType || lex.Value != lts.lexemes[i].Value || lts.lexemes[i].Pos != lex.Pos {
				t.Errorf("Test %q, lexeme %d: Expected %#v. Got %#v instead", lts.name, i, lts.lexemes[i], lex)
			}
		}
	}
}
