package lingo

import (
	"fmt"
	"strings"
)

// DependencyType represents the relation between two words
type DependencyType byte

var dependencyTypeLookup map[string]DependencyType

func init() {
	dependencyTypeLookup = make(map[string]DependencyType)
	for dt := NoDepType; dt < MAXDEPTYPE; dt++ {
		s := dt.String()
		dependencyTypeLookup[s] = DependencyType(dt)
		dependencyTypeLookup[strings.ToLower(s)] = DependencyType(dt)
	}
}

func (dt DependencyType) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("%v", dt)), nil
}

func (dt *DependencyType) UnmarshalText(text []byte) error {
	str := strings.Trim(string(text), `"`) // for JSON use, if any
	deptype, _ := dependencyTypeLookup[str]
	*dt = deptype
	return nil
}

// list of dependency type functions

func InDepTypes(x DependencyType, set []DependencyType) bool {
	for _, v := range set {
		if v == x {
			return true
		}
	}
	return false
}

func IsModifier(x DependencyType) bool      { return InDepTypes(x, Modifiers) }
func IsCompound(x DependencyType) bool      { return InDepTypes(x, Compounds) }
func IsDeterminerRel(x DependencyType) bool { return InDepTypes(x, DeterminerRels) }
func IsMultiword(x DependencyType) bool     { return InDepTypes(x, MultiWord) }
func IsQuantifier(x DependencyType) bool    { return InDepTypes(x, QuantifingMods) }
