package corpus

import (
	"regexp"

	"github.com/chewxy/lingo"
)

type conversionPattern struct {
	pattern     *regexp.Regexp
	replacement string
}

func newConversionPattern(from, to string) conversionPattern {
	rFrom := regexp.MustCompile(from)
	return conversionPattern{rFrom, to}
}

// plural -> singular
var plural = []conversionPattern{
	newConversionPattern("(quiz)$", "${1}zes"),
	newConversionPattern("^(ox)$", "${1}en"),
	newConversionPattern("([m|l])ouse$", "${1}ice"),
	newConversionPattern("(matr|vert|ind)ix|ex$", "${1}ices"),
	newConversionPattern("(x|ch|ss|sh)$", "${1}es"),
	newConversionPattern("([^aeiouy]|qu)ies$", "${1}y"),
	newConversionPattern("([^aeiouy]|qu)y$", "${1}ies"),
	newConversionPattern("(hive)$", "${1}s"),
	newConversionPattern("(?:([^f])fe|([lr])f)$", "${1}${2}ves"),
	newConversionPattern("sis$", "ses"),
	newConversionPattern("([ti])um$", "${1}a"),
	newConversionPattern("(buffal|tomat|potat)o$", "${1}oes"),
	newConversionPattern("(bu)s$", "${1}ses"),
	newConversionPattern("(alias|status|sex)$", "${1}es"),
	newConversionPattern("(octop|vir)us$", "${1}i"),
	newConversionPattern("(ax|test)is$", "${1}es"),
	newConversionPattern("s$", "s"),
	newConversionPattern("$", "s"),
}

// singular -> plural
var singular = []conversionPattern{
	newConversionPattern("(quiz)zes$", "${1}"),
	newConversionPattern("(matr)ices$", "${1}ix"),
	newConversionPattern("(vert|ind)ices$", "${1}ex"),
	newConversionPattern("^(ox)en", "${1}"),
	newConversionPattern("(alias|status)es$", "${1}"),
	newConversionPattern("(octop|vir)i$", "${1}us"),
	newConversionPattern("(cris|ax|test)es$", "${1}is"),
	newConversionPattern("(shoe)s$", "${1}"),
	newConversionPattern("(o)es$", "${1}"),
	newConversionPattern("(bus)es$", "${1}"),
	newConversionPattern("([m|l])ice$", "${1}ouse"),
	newConversionPattern("(x|ch|ss|sh)es$", "${1}"),
	newConversionPattern("(m)ovies$", "${1}ovie"),
	newConversionPattern("(s)eries$", "${1}eries"),
	newConversionPattern("([^aeiouy]|qu)ies$", "${1}y"),
	newConversionPattern("([lr])ves$", "${1}f"),
	newConversionPattern("(tive)s$", "${1}"),
	newConversionPattern("(hive)s$", "${1}"),
	newConversionPattern("([^f])ves$", "${1}fe"),
	newConversionPattern("(^analy)ses$", "${1}sis"),
	newConversionPattern("((a)naly|(b)a|(d)iagno|(p)arenthe|(p)rogno|(s)ynop|(t)he)ses$", "${1}${2}sis"),
	newConversionPattern("([ti])a$", "${1}um"),
	newConversionPattern("(n)ews$", "${1}ews"),
	newConversionPattern("s$", ""),
}

// weird pluralizations that don't match the rules above
var irregular = []conversionPattern{
	newConversionPattern("person", "people"),
	newConversionPattern("man", "men"),
	newConversionPattern("child", "children"),
	newConversionPattern("sex", "sexes"),
	newConversionPattern("move", "moves"),
	newConversionPattern("sleeve", "sleeves"),
	newConversionPattern("datum", "data"),
	newConversionPattern("box", "boxes"),
	newConversionPattern("knife", "knives"),
}

var unconvertable = []string{
	"equipment",
	"information",
	"rice",
	"money",
	"species",
	"series",
	"fish",
	"sheep",
}

// Pluralize pluralizes words based on rules known
func Pluralize(word string) string {
	if lingo.InStringSlice(word, unconvertable) {
		return word
	}

	for _, cp := range irregular {
		if cp.pattern.MatchString(word) {
			return cp.replacement
		}
	}

	for _, cp := range plural {
		if cp.pattern.MatchString(word) {
			// log.Printf("\t%q Matches %q", word, cp.pattern.String())
			return cp.pattern.ReplaceAllString(word, cp.replacement)
		}
	}
	return word
}

// Singularize singularizes words based on rules known
func Singularize(word string) string {
	if lingo.InStringSlice(word, unconvertable) {
		return word
	}

	for _, cp := range singular {
		if cp.pattern.MatchString(word) {
			return cp.pattern.ReplaceAllString(word, cp.replacement)
		}
	}
	return word
}
