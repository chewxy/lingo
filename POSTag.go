package lingo

import (
	"fmt"
	"strings"
)

// POSTag represents a Part of Speech Tag.
type POSTag byte

var posTagLookup map[string]POSTag

func init() {
	posTagLookup = make(map[string]POSTag)
	for t := X; t < MAXTAG; t++ {
		s := t.String()
		posTagLookup[s] = POSTag(t)
		posTagLookup[strings.ToLower(s)] = POSTag(t)
	}
}

func (p POSTag) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("%v", p)), nil // add quotes back
}

func (p *POSTag) UnmarshalText(text []byte) error {
	str := strings.Trim(string(text), `"`) // for JSON use, if any
	tag, _ := posTagLookup[str]
	*p = tag
	return nil
}

// POSTag related functions
func InPOSTags(x POSTag, set []POSTag) bool {
	for _, v := range set {
		if v == x {
			return true
		}
	}
	return false
}

func IsAdjective(x POSTag) bool     { return InPOSTags(x, Adjectives) }
func IsNoun(x POSTag) bool          { return InPOSTags(x, Nouns) }
func IsProperNoun(x POSTag) bool    { return InPOSTags(x, ProperNouns) }
func IsVerb(x POSTag) bool          { return InPOSTags(x, Verbs) }
func IsAdverb(x POSTag) bool        { return InPOSTags(x, Adverbs) }
func IsInterrogative(x POSTag) bool { return InPOSTags(x, Interrogatives) }
func IsDeterminer(x POSTag) bool    { return InPOSTags(x, Determiners) }
func IsNumber(x POSTag) bool        { return InPOSTags(x, Numbers) }
func IsSymbol(x POSTag) bool        { return InPOSTags(x, Symbols) }
