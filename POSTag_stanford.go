// +build stanfordtags

package lingo

//go:generate stringer -type=POSTag -output=POSTag_stanford_string.go

const BUILD_TAGSET = "stanfordtags"

const (
	X           POSTag = iota // aka NULLTAG
	UNKNOWN_TAG               // Unknown
	ROOT_TAG                  // For Root
	CC                        // Coordinating conjunction
	CD                        // Cardinal number
	DT                        // Determiner
	EX                        // Existential there
	FW                        // Foreign word
	IN                        // Preposition or subordinating conjunction
	JJ                        // Adjective
	JJR                       // Adjective, comparative
	JJS                       // Adjective, superlative
	LS                        // List item marker
	MD                        // Modal
	NN                        // Noun, singular or mass
	NNS                       // Noun, plural
	NNP                       // Proper noun, singular
	NNPS                      // Proper noun, plural
	PDT                       // Predeterminer
	POS                       // Possessive ending
	PRP                       // Personal pronoun
	PPRP                      // Possessive pronoun (PRP$)
	RB                        // Adverb
	RBR                       // Adverb, comparative
	RBS                       // Adverb, superlative
	RP                        // Particle
	SYM                       // Symbol
	TO                        // to
	UH                        // Interjection
	VB                        // Verb, base form
	VBD                       // Verb, past tense
	VBG                       // Verb, gerund or present participle
	VBN                       // Verb, past participle
	VBP                       // Verb, non-3rd person singular present
	VBZ                       // Verb, 3rd person singular present
	WDT                       // Wh-determiner
	WP                        // Wh-pronoun
	PWP                       // Possessive wh-pronoun (WP$)
	WRB                       // Wh-adverb

	// Punctuation related stuff: http://stackoverflow.com/a/21546294
	COMMA      // Obvious isn't it?
	FULLSTOP   // fullstop
	OPENQUOTE  // Penn Treebank uses ``
	CLOSEQUOTE // Penn Treebank uses ''
	COLON
	DOLLAR
	HASHSIGN
	LEFTBRACE
	RIGHTBRACE

	// Extensions for web shit: https://www.ldc.upenn.edu/sites/www.ldc.upenn.edu/files/etb-supplementary-guidelines-2009-addendum.pdf
	// http://clear.colorado.edu/compsem/documents/treebank_guidelines.pdf
	HYPH // Hyphen in split compounds
	AFX  // affix
	ADD  // url or email addy
	NFP  // superfluous (non final) puncutation
	GW   // Goes WIth
	XX   // deidentified data (aka giberish)

	MAXTAG
)

// POSTagShortcut is a shortcut function to help the POSTagger shortcircuit some decisions about what the tag is
func POSTagShortcut(l Lexeme) (POSTag, bool) {
	switch l.LexemeType {
	case Number:
		return CD, true
	case Punctuation:
		switch l.Value {
		case ",":
			return COMMA, true
		case ".":
			return FULLSTOP, true
		case "``":
			return OPENQUOTE, true
		case "''":
			return CLOSEQUOTE, true
		case ":":
			return COLON, true
		case "#":
			return HASHSIGN, true
		case "(":
			return LEFTBRACE, true
		case ")":
			return RIGHTBRACE, true
		default:
			return X, false
		}
	case Symbol:
		return SYM, true
	case URI:
		return ADD, true
	case Date:
		return CD, true
	case Time:
		return CD, true
	case EOF:
		return X, true
	}
	return X, false
}

// sets

var Adjectives = []POSTag{JJ, JJR, JJS}
var Nouns = []POSTag{NN, NNP, NNS, NNPS}
var ProperNouns = []POSTag{NNP, NNPS}
var Verbs = []POSTag{VB, VBD, VBG, VBN, VBP, VBZ}
var Adverbs = []POSTag{RB, RBR, RBS}
var Determiners = []POSTag{DT, PDT}
var Interrogatives = []POSTag{WDT, WP, PWP, WRB}
var Numbers = []POSTag{CD}
var Symbols = []POSTag{SYM, FULLSTOP, COMMA, OPENQUOTE, COLON, DOLLAR, HASHSIGN, LEFTBRACE, RIGHTBRACE, HYPH, NFP}

// IsIN returns true if the POSTag is a subordinating conjunction.
// The reason why this exists is because in the stanford tag, IN is the POSTag
// while in the universal dependencies, it's the SCONJ POSTag
func IsIN(x POSTag) bool { return x == IN }
