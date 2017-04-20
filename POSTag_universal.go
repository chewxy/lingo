// +build !stanfordtags

package lingo

//go:generate stringer -type=POSTag -output=POSTag_universal_string.go

const BUILD_TAGSET = "universaltags"

const (
	X POSTag = iota // aka NULLTAG
	UNKNOWN_TAG
	ROOT_TAG
	ADJ
	ADP
	ADV
	AUX
	CONJ
	DET
	INTJ
	NOUN
	NUM
	PART
	PRON
	PROPN
	PUNCT
	SCONJ
	SYM
	VERB

	MAXTAG // MAXTAG is provided here as index support
)

// POSTagShortcut is a shortcut function to help the POSTagger shortcircuit some decisions about what the tag is
func POSTagShortcut(l Lexeme) (POSTag, bool) {
	switch l.LexemeType {
	case Number:
		return NUM, true
	case Punctuation:
		return PUNCT, true
	case Symbol:
		return SYM, true
	case URI:
		return X, true
	case Date:
		return NUM, true
	case Time:
		return NUM, true
	case EOF:
		return X, true
	}
	return X, false
}

var Adjectives = []POSTag{ADJ}
var Nouns = []POSTag{NOUN, PROPN}
var ProperNouns = []POSTag{PROPN}
var Verbs = []POSTag{VERB}
var Adverbs = []POSTag{ADV}
var Determiners = []POSTag{DET}
var Interrogatives = []POSTag{PRON, DET, ADV}
var Numbers = []POSTag{NUM}
var Symbols = []POSTag{SYM, PUNCT}

// IsIN returns true if the POSTag is a subordinating conjunction.
// The reason why this exists is because in the stanford tag, IN is the POSTag
// while in the universal dependencies, it's the SCONJ POSTag
func IsIN(x POSTag) bool { return x == SCONJ }
