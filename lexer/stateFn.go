package lexer

import (
	"unicode"

	"github.com/chewxy/lingo"
)

type stateFn func(*Lexer) stateFn

func lexText(l *Lexer) (fn stateFn) {
	for {
		next := l.next()
		if next == eof {
			break
		}
		if l.pos != l.start {
			switch {
			case unicode.IsSpace(next):
				l.backup()
				fn = lexWhitespace
			case unicode.IsDigit(next):

				// if the position is start +1.
				// This means that the first char of the string to be lexed is a number
				// this prevents things like "yay1111" to be lexed as "yay" and "1111"
				if l.pos == l.start+1 {
					l.backup()
					return lexNumber
				}
			case next == ':':
				if l.peek() == '/' {
					l.accept() // accept ':'
					l.next()
					if l.peek() == '/' {
						l.accept()
						return lexURI
					}
					// otherwise...
					l.backup()
					// "unaccept". since '/' has a width of 1 we can do the following
					l.buf.Truncate(l.buf.Len() - 1)
				}
				fn = lexPunctuation
			case unicode.IsPunct(next):
				// for things like "u.s" or "i.e." or "e.g."
				n := l.peek()

				switch {
				case next == '\'':
					if unicode.IsLetter(n) {
						l.emit(lingo.Word)
						return lexText
					}
				case n == eof:
					// common scenario - where a punctuation ends the sentence, and this thing is unable to backup
					l.width = 1
					l.backup()
					l.width = 0
					fn = lexPunctuation
					goto finishup // goto because there are other cases below
				case unicode.IsLetter(n):
					l.accept()
					return lexText
				default:
					// it's definitely a punctuation
					l.backup()
					fn = lexPunctuation
				}

			case unicode.IsSymbol(next):
				// for things like "e-mail"
				n := l.next()
				if unicode.IsLetter(n) {
					return lexText
				}

				l.backup()
				fn = lexSymbol
			case next == 'n':
				// for things like "don't" or "doesn't"
				n := l.peek()
				if n == '\'' {
					l.backup()
					l.emit(lingo.Word)
					return lexPunctuation
				} else {
					l.accept() // accept n
					return lexText
				}
			}
		}

	finishup:
		if fn != nil {
			if l.start != l.pos {
				l.emit(lingo.Word)
			}
			return fn
		}
		// otherwise keep lexText
		l.accept()
	}

	if l.pos > l.start {
		l.emit(lingo.Word)
	}

	l.emit(lingo.EOF)
	return nil
}

// lexNumber lexes numbers. It accepts runs of unicode digits.
// Upon stopping, it checks to see if the next value is a '.'. If it is, then it's a decimal value, and continues a run
// Upon stopping a second time, it checks for 'e' or 'E', for exponentiation - 1.2E2
func lexNumber(l *Lexer) (fn stateFn) {
	l.acceptRunFn(unicode.IsDigit)

	next := l.next()
	switch next {
	case '.':
		l.accept() // accept the dot
		l.acceptRunFn(unicode.IsDigit)
	case '-', '/':
		// standardize
		l.r = '/'
		l.accept()
		return lexDate
	case ':':
		if l.pos-l.start == 3 {
			l.accept()
			return lexTime
		} else {
			l.backup()
			l.emit(lingo.Number)
			return lexPunctuation
		}
	default:
		l.backup()
	}

	if l.acceptRun("eE") {
		// handle negative exponents
		if l.peek() == '-' {
			l.next()
			l.accept()
			return lexNumber(l)
		}
		l.acceptRunFn(unicode.IsDigit)
	}
	l.backup()

	if l.buf.Len() == 1 && l.buf.Bytes()[0] == '-' {
		l.emit(lingo.Punctuation) // dash
		return lexWhitespace
	}
	l.emit(lingo.Number)
	return lexWhitespace
}

func lexWhitespace(l *Lexer) (fn stateFn) {
	l.acceptRunFn(unicode.IsSpace)
	l.lineCount()
	// l.incrementLineCount()
	// l.backup()
	l.ignore() //nothing will be emitted

	next := l.peek()
	switch {
	case unicode.IsDigit(next):
		return lexNumber
	case unicode.IsPunct(next):
		if next == '-' {
			l.next()
			l.accept()
			return lexNumber
		}
		return lexPunctuation
	case unicode.IsSymbol(next):
		return lexSymbol
	}

	return lexText
}

func lexPunctuation(l *Lexer) (fn stateFn) {
	next := l.next()

	switch next {
	case '\'':
		l.accept()
		n := l.peek()
		switch n {
		case 't', 's', 'm', 'd':
			l.next()
			l.accept() // accept 't'/'s'...
			l.emit(lingo.Word)
			return lexWhitespace
		}
	case '.':
		l.accept()
		// for cases such as "U.S" or "i.e"
		n := l.peek()
		if unicode.IsLetter(n) {
			l.accept() // accept .
			l.next()
			l.accept()
			return lexText
		}
	default:
		// log.Printf("Next %q", next)
	}

	accepted := l.acceptRunFn(unicode.IsPunct) // check for any other runs of punctuations
	if accepted == 0 && unicode.IsPunct(next) {
		l.accept()
	}
	l.emit(lingo.Punctuation)
	return lexWhitespace
}

func lexSymbol(l *Lexer) (fn stateFn) {
	l.acceptRunFn(unicode.IsSymbol)
	l.acceptRunFn(unicode.IsPunct) // any symbol punctuation combination should be treated as a symbole
	l.emit(lingo.Symbol)
	return lexWhitespace
}

func lexURI(l *Lexer) (fn stateFn) {
	eof := l.nextUntilEOF(" ")
	if !eof {
		l.backup()
		l.backup()
		next := l.next()
		if unicode.IsPunct(next) {
			l.backup()
			l.emit(lingo.URI)
			return lexPunctuation
		}
	}

	l.emit(lingo.URI)
	return lexWhitespace
}

func lexDate(l *Lexer) (fn stateFn) {
	l.acceptRunFn(unicode.IsDigit)
	next := l.next()
	if next != '/' && next != '-' {
		l.backup()
		l.emit(lingo.Number) // fractions are numbers
		return lexWhitespace
	}
	l.r = '/' // standardize
	l.accept()

	l.acceptRunFn(unicode.IsDigit)
	l.emit(lingo.Date)
	return lexWhitespace
}

func lexTime(l *Lexer) (fn stateFn) {
	l.acceptRunFn(unicode.IsDigit)
	next := l.next()
	if next != ':' {
		l.backup()
		l.emit(lingo.Time)
		return lexWhitespace
	}
	l.accept()
	l.acceptRunFn(unicode.IsDigit)
	l.emit(lingo.Time)
	return lexWhitespace
}
