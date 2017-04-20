// +build stanfordrel

package lingo

const BUILD_RELSET = "stanfordrel"

//go:generate stringer -type=DependencyType -output=dependencyType_stanford_string.go

// http://nlp.stanford.edu/software/dependencies_manual.pdf
const (
	NoDepType DependencyType = iota
	Dep
	Root
	Aux           // Auxilliary
	AuxPass       // passive auxiliary
	Cop           // Copula
	Arg           // argument
	Agent         // agent
	Comp          // Complement
	AComp         // adjectival complement
	CComp         // clausal complement with internal subject
	XComp         // clausal complement with external subject
	Obj           // Object
	DObj          // Direct Object
	IObj          // Indirect Object
	PObj          // Object of preposition
	Subj          // subject
	NSubj         // Nominal Subject
	NSubjPass     // passive nominal subject
	CSubj         // clausal subject
	CSubjPass     // passive clausal subject
	Coordination  // coordination (cannot use CC, as it's a POSTag)
	Conj          // conjunction
	Expl          // Expletive
	Mod           // modifier
	AMod          // adjectival modifier
	Appos         // Appositional modifier
	Advcl         // adverbial clause modifier
	Det           // determiner
	Predet        // predeterminer
	Preconj       // Preconjunction
	Vmod          // reduced, nonfinite verbal modifier
	MWE           // multiword expression modifier
	Mark          // marker (word introducing an Advcl or CComp)
	AdvMod        // adverbial modifier
	Neg           // negation modifier
	RCMod         // relative clause modifier
	QuantMod      // quantifier modifier
	NounMod       // Noun Compound Modifier (cannot use NN because NN is defined as a POSTag)
	NPAdvMod      // Noun phrase adverbial modifier
	TMod          // temporal modifier
	Num           // Numeric Modifier
	NumberElement // element of compound number (cannot use Number because Number is defined as a LexemeType)
	Prep          // prepositional modifier
	Poss          // possession modifier
	Possessive    // possessive modifier ('s)
	PRT           // phrasal verb partical
	Parataxis     // Parataxis (words that are next to each other)
	GoesWith      // GoesWith
	Punct         // punctuation
	Ref           // referant
	SDep          // Semantic Dependent
	XSubj         // controlling subject

	// additional stuff not found in the original, but found in EWT
	Case
	Compound
	NMod
	Discourse
	NumMod
	RelCl
	NFinCl
	NMod_Poss
	NMod_NPMod
	Vocative
	List
	MWPrep // multiword prepositional modifier
	Remnant
	Acl
	NPMod
	MDVod
	DetMod

	// found in stanford nnparser SD models
	PComp

	MAXDEPTYPE
)

var Modifiers = []DependencyType{AMod}
var Compounds = []DependencyType{Compound}
var DeterminerRels = []DependencyType{Det, DetMod}
var MultiWord = []DependencyType{MWE, MWPrep, Compound, Parataxis}
var QuantifingMods = []DependencyType{QuantMod, NumMod}
