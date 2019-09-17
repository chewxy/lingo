package lexer

import (
	"bufio"
	"bytes"
	"io"
	"strings"

	"golang.org/x/text/unicode/norm"

	"github.com/chewxy/lingo"
)

const eof rune = -1

type Lexer struct {
	name  string
	input *bufio.Reader

	state stateFn
	r     rune
	width int
	pos   int
	start int
	line  int
	col   int

	// the string we're reading
	buf *bytes.Buffer

	Output chan lingo.Lexeme
	Errors chan error
}

func New(name string, r io.Reader) *Lexer {
	return &Lexer{
		name:  name,
		input: bufio.NewReader(r),

		width: 1,
		start: 1, // for userfriendliness, the column index starts at 1
		col:   1,
		pos:   1,
		buf:   new(bytes.Buffer),

		Output: make(chan lingo.Lexeme),
		Errors: make(chan error),
	}
}

func (l *Lexer) Run() {
	defer close(l.Output)
	for state := lexText; state != nil; {
		state = state(l)
	}
}

// Reset resets the buffers. It creates a new Output and Error channel
func (l *Lexer) Reset(r io.Reader) {
	l.input.Reset(r)
	l.buf.Reset()
	l.Output = make(chan lingo.Lexeme)
	l.Errors = make(chan error)
}

func (l *Lexer) next() rune {
	var err error
	l.r, l.width, err = l.input.ReadRune()
	if err == io.EOF {
		l.width = 1
		return eof
	}
	l.col += l.width
	l.pos += l.width

	return l.r
}

// nextUntilEOF will loop until it finds the matching string OR EOF
func (l *Lexer) nextUntilEOF(s string) bool {
	for r := l.next(); r != eof && strings.IndexRune(s, r) < 0; r = l.next() {
		// l.next()
		l.accept()
	}
	if l.r == eof {
		return true
	}
	return false
}

func (l *Lexer) backup() {
	l.input.UnreadRune()
	l.pos -= l.width
	l.col -= l.width
}

func (l *Lexer) peek() rune {
	backup := l.r
	pos := l.pos
	col := l.col

	r := l.next()
	l.backup()

	l.pos = pos
	l.col = col
	l.r = backup
	return r
}

func (l *Lexer) lineCount() {
	newLines := bytes.Count(l.buf.Bytes(), []byte("\n"))

	l.line += newLines
	if newLines > 0 {
		l.col = 1
	}
}

func (l *Lexer) accept() {
	l.buf.WriteRune(l.r)
}

func (l *Lexer) acceptRun(valid string) (accepted bool) {
	for strings.IndexRune(valid, l.peek()) >= 0 {
		l.next()
		l.accept()
		accepted = true
	}
	return
}

func (l *Lexer) acceptRunFn(fn func(rune) bool) (accepted int) {
	for fn(l.peek()) {
		l.next()
		l.accept()
		accepted++
	}
	return
}

func (l *Lexer) ignore() {
	l.start = l.pos
	l.buf.Reset()
}

func (l *Lexer) emit(t lingo.LexemeType) {
	normalized := string(norm.NFC.Bytes(l.buf.Bytes()))
	lex := lingo.MakeLexeme(normalized, t)
	lex.Line = l.line
	lex.Col = l.start

	l.Output <- lex

	// reset
	l.ignore()
	if l.r != 0x0 {
		l.buf.WriteRune(l.r)
	}
}
