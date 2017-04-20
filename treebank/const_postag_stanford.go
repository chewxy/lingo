// +build stanfordtags

package treebank

import "github.com/chewxy/lingo"

var posTagTable map[string]lingo.POSTag = map[string]lingo.POSTag{
	"X": lingo.X,

	"CC":   lingo.CC,
	"CD":   lingo.CD,
	"DT":   lingo.DT,
	"EX":   lingo.EX,
	"FW":   lingo.FW,
	"IN":   lingo.IN,
	"JJ":   lingo.JJ,
	"JJR":  lingo.JJR,
	"JJS":  lingo.JJS,
	"LS":   lingo.LS,
	"MD":   lingo.MD,
	"NN":   lingo.NN,
	"NNS":  lingo.NNS,
	"NNP":  lingo.NNP,
	"NNPS": lingo.NNPS,
	"PDT":  lingo.PDT,
	"POS":  lingo.POS,
	"PRP":  lingo.PRP,
	"PPRP": lingo.PPRP,
	"PRP$": lingo.PPRP,
	"RB":   lingo.RB,
	"RBR":  lingo.RBR,
	"RBS":  lingo.RBS,
	"RP":   lingo.RP,
	"SYM":  lingo.SYM,
	"TO":   lingo.TO,
	"UH":   lingo.UH,
	"VB":   lingo.VB,
	"VBD":  lingo.VBD,
	"VBG":  lingo.VBG,
	"VBN":  lingo.VBN,
	"VBP":  lingo.VBP,
	"VBZ":  lingo.VBZ,
	"WDT":  lingo.WDT,
	"WP":   lingo.WP,
	"PWP":  lingo.PWP,
	"WP$":  lingo.PWP,
	"WRB":  lingo.WRB,

	// punctuation
	",":     lingo.COMMA,
	"``":    lingo.OPENQUOTE,
	"''":    lingo.CLOSEQUOTE,
	".":     lingo.FULLSTOP,
	":":     lingo.COLON,
	"$":     lingo.DOLLAR,
	"#":     lingo.HASHSIGN,
	"-LRB-": lingo.LEFTBRACE,
	"-RRB-": lingo.RIGHTBRACE,

	"ADD":  lingo.ADD,
	"NFP":  lingo.NFP,
	"HYPH": lingo.HYPH,
	"GW":   lingo.GW,
	"AFX":  lingo.AFX,
	"XX":   lingo.XX,

	"-NULL-":    lingo.X,
	"-ROOT-":    lingo.ROOT_TAG,
	"-UNKNOWN-": lingo.UNKNOWN_TAG,
}
