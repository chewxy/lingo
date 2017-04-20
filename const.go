package lingo

// constants that are not pertaining to build tags

var empty struct{}

// NumberWords was generated with this python code
/*
	numberWords = {}

	simple = '''zero one two three four five six seven eight nine ten eleven twelve
	        thirteen fourteen fifteen sixteen seventeen eighteen nineteen
	        twenty'''.split()
	for i, word in zip(xrange(0, 20+1), simple):
	    numberWords[word] = i

	tense = '''thirty forty fifty sixty seventy eighty ninety hundred'''.split()
	for i, word in zip(xrange(30, 100+1, 10), tense):
		numberWords[word] = i

	larges = '''thousand million billion trillion quadrillion quintillion sextillion septillion'''.split()
	for i, word in zip(xrange(3, 24+1, 3), larges):
		numberWords[word] = 10**i
*/
var NumberWords = map[string]int{
	"zero":        0,
	"one":         1,
	"two":         2,
	"three":       3,
	"four":        4,
	"five":        5,
	"six":         6,
	"seven":       7,
	"eight":       8,
	"nine":        9,
	"ten":         10,
	"eleven":      11,
	"twelve":      12,
	"thirteen":    13,
	"fourteen":    14,
	"fifteen":     15,
	"sixteen":     16,
	"nineteen":    19,
	"seventeen":   17,
	"eighteen":    18,
	"twenty":      20,
	"thirty":      30,
	"forty":       40,
	"fifty":       50,
	"sixty":       60,
	"seventy":     70,
	"eighty":      80,
	"ninety":      90,
	"hundred":     100,
	"thousand":    1000,
	"million":     1000000,
	"billion":     1000000000,
	"trillion":    1000000000000,
	"quadrillion": 1000000000000000,
	// "quintillion": 1000000000000000000,
	// "sextillion": 1000000000000000000000,
	// "septillion": 1000000000000000000000000,
}
