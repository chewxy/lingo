package lingo

import (
	"bytes"
	"unicode"
)

// Shape represents the shape of a word. It's currently implemented as a string
type Shape string

func (l Lexeme) Shape() Shape {
	s := l.Value

	if len(s) > 50 {
		return Shape("Long")
	}

	var buf bytes.Buffer

	previousCharShape := ' '
	currentCharShape := ' '
	sequence := 0
	for _, c := range s {
		switch {
		case unicode.IsLetter(c):
			if unicode.IsUpper(c) {
				currentCharShape = 'X'
			} else {
				currentCharShape = 'x'
			}

		case unicode.IsDigit(c):
			currentCharShape = 'd'

		case l.LexemeType == URI:
			return Shape("URI")

		default:
			currentCharShape = c
		}

		if previousCharShape == currentCharShape {
			sequence++
		} else {
			sequence = 0 // reset the sequence
			previousCharShape = currentCharShape
		}

		if sequence < 4 {
			buf.WriteRune(currentCharShape)
		}
	}

	retVal := buf.String()

	return Shape(retVal)
}
