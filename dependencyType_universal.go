// +build !stanfordrel

package lingo

const BUILD_RELSET = "universalrel"

//go:generate stringer -type=DependencyType -output=dependencyType_universal_string.go

// http://universaldependencies.github.io/docs/en/dep/all.html
const (
	NoDepType DependencyType = iota
	Dep
	Root

	// Core dependents of clausal predicates

	// nominal dependencies
	NSubj
	NSubjPass
	DObj
	IObj

	// predicate dependencies
	CSubj
	CSubjPass
	CComp

	XComp

	// Noun dependents

	// nominal dependencies
	NumMod
	Appos
	NMod

	// predicate dependencies
	ACl
	ACl_RelCl // RCMod in stanford deps
	Det
	Det_PreDet

	// modifier word
	AMod
	Neg

	// Case Marking, preposition, possessive
	Case

	//Non-Core Dependents of Clausal Predicates

	// Nominal dependencies
	NMod_NPMod
	NMod_TMod
	NMod_Poss

	// Predicate Dependencies
	AdvCl

	// Modifier Word
	AdvMod

	// Compounding and Unanalyzed
	Compound
	Compound_Part
	Name // Unused in English
	MWE
	Foreign  // Unused in English
	GoesWith // Unused in English

	// Loose Joining Relations
	List
	Dislocated // Unused in English
	Parataxis
	Remnant    // Unused in English
	Reparandum // Unused in English

	// Special Clausal Dependents

	// Nominal Dependent
	Vocative // Unused in English
	Discourse
	Expl

	// Auxilliary
	Aux
	AuxPass
	Cop

	// Other
	Mark
	Punct

	// Coordination

	Conj
	Coordination // CC
	CC_PreConj

	MAXDEPTYPE
)

var Modifiers = []DependencyType{AMod}
var Compounds = []DependencyType{Compound, Compound_Part}
var DeterminerRels = []DependencyType{Det, Det_PreDet}
var MultiWord = []DependencyType{MWE, Compound, Compound_Part, Parataxis}
var QuantifingMods = []DependencyType{NumMod}
