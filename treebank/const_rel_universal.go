// +build !stanfordrel

package treebank

import "github.com/chewxy/lingo"

var dependencyTable map[string]lingo.DependencyType = map[string]lingo.DependencyType{
	"dep":          lingo.Dep,
	"root":         lingo.Root,
	"nsubj":        lingo.NSubj,
	"nsubjpass":    lingo.NSubjPass,
	"dobj":         lingo.DObj,
	"iobj":         lingo.IObj,
	"csubj":        lingo.CSubj,
	"csubjpass":    lingo.CSubjPass,
	"ccomp":        lingo.CComp,
	"xcomp":        lingo.XComp,
	"nummod":       lingo.NumMod,
	"appos":        lingo.Appos,
	"nmod":         lingo.NMod,
	"acl":          lingo.ACl,
	"acl:relcl":    lingo.ACl_RelCl,
	"det":          lingo.Det,
	"det:predet":   lingo.Det_PreDet,
	"amod":         lingo.AMod,
	"neg":          lingo.Neg,
	"case":         lingo.Case,
	"nmod:npmod":   lingo.NMod_NPMod,
	"nmod:tmod":    lingo.NMod_TMod,
	"nmod:poss":    lingo.NMod_Poss,
	"advcl":        lingo.AdvCl,
	"advmod":       lingo.AdvMod,
	"compound":     lingo.Compound,
	"compound:prt": lingo.Compound_Part,
	"name":         lingo.Name,
	"mwe":          lingo.MWE,
	"foreign":      lingo.Foreign,
	"goeswith":     lingo.GoesWith,
	"list":         lingo.List,
	"dislocated":   lingo.Dislocated,
	"parataxis":    lingo.Parataxis,
	"remnant":      lingo.Remnant,
	"reparandum":   lingo.Reparandum,
	"vocative":     lingo.Vocative,
	"discourse":    lingo.Discourse,
	"expl":         lingo.Expl,
	"aux":          lingo.Aux,
	"auxpass":      lingo.AuxPass,
	"cop":          lingo.Cop,
	"mark":         lingo.Mark,
	"punct":        lingo.Punct,
	"conj":         lingo.Conj,
	"cc":           lingo.Coordination,
	"cc:preconj":   lingo.CC_PreConj, // https://github.com/UniversalDependencies/docs/issues/221
	"conj:preconj": lingo.CC_PreConj, // https://github.com/UniversalDependencies/docs/issues/221

	"-NULL-": lingo.NoDepType,
}
