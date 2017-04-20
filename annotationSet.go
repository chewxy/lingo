package lingo

import (
	"sort"
	"unsafe"

	"github.com/xtgo/set"
)

type AnnotationSet []*Annotation

func (as AnnotationSet) Len() int      { return len(as) }
func (as AnnotationSet) Swap(i, j int) { as[i], as[j] = as[j], as[i] }
func (as AnnotationSet) Less(i, j int) bool {
	return uintptr(unsafe.Pointer(as[i])) < uintptr(unsafe.Pointer(as[j]))
}

func (as AnnotationSet) Set() AnnotationSet {
	sort.Sort(as)
	n := set.Uniq(as)
	return as[:n]
}

func (as AnnotationSet) Contains(a *Annotation) bool {
	if as.Index(a) == len(as) {
		return false
	}
	return true
}

func (as AnnotationSet) Index(a *Annotation) int {
	for i, an := range as {
		if an == a {
			return i
		}
	}
	return len(as)
}

func (as AnnotationSet) Add(a *Annotation) AnnotationSet {
	if as.Contains(a) {
		return as
	}
	as = append(as, a)
	return as
}
