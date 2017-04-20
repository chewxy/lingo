package dep

import (
	"log"

	"github.com/chewxy/lingo"
)

// applies common fixes
func fix(d *lingo.Dependency) {
	// NNP fix:
	// If a sentence is [a, b, c, D, E, f, g]
	// where D, E are NNPs, they should be compound words
	// The head should be the one with higher headID
	spans := properNounSpans(d)
	for _, s := range spans {
		// we don't care about single word proper nouns
		if s.end-s.start <= 1 {
			continue
		}

		phrase := d.AnnotatedSentence[s.start:s.end]

		// pick up all compound roots
		// find annotations that do not have compound as deptype
		var compoundRoots lingo.AnnotationSet
		var problematic lingo.AnnotationSet
		for _, a := range phrase {
			if lingo.IsCompound(a.DependencyType) {
				compoundRoots = compoundRoots.Add(a.Head)
			}

			if !lingo.IsCompound(a.DependencyType) && a.ID != s.end-1 {
				problematic = problematic.Add(a)
			}
		}

		// if no root
		if len(compoundRoots) == 0 {
			// actual root is the word with the largest ID
			var compoundRoot *lingo.Annotation
			var rootRoot *lingo.Annotation
			for last := -1; s.end+last >= s.start; last-- {
				predictedRoot := s.end + last
				compoundRoot = d.AnnotatedSentence[predictedRoot]

				// incorrects :
				//	dep==Dep
				// 	dep==Root && others has dep != root

				if compoundRoot.DependencyType == lingo.Dep {
					problematic = problematic.Add(compoundRoot)
					continue
				}

				if compoundRoot.DependencyType != lingo.Dep && compoundRoot.DependencyType != lingo.Root {
					break
				}

				if compoundRoot.DependencyType == lingo.Root {
					rootRoot = compoundRoot
					problematic = problematic.Add(compoundRoot)
				}
			}

			if rootRoot != nil && rootRoot != compoundRoot {
				// we have two potential roots. Choose the best
				log.Println("Problem when fixing: more than one possible compound root found")
			}

			for _, a := range problematic {
				if a == compoundRoot {
					continue
				}
				tmpHead := a.Head
				tmpRel := a.DependencyType

				a.SetHead(compoundRoot)
				a.DependencyType = lingo.Compound

				for _, childID := range d.AnnotatedSentence.Children(a.ID) {
					childA := d.AnnotatedSentence[childID]
					childA.SetHead(tmpHead)
					childA.DependencyType = tmpRel
				}
			}

		}

		// if more than one root...
		logf("More than zero compound roots not handled yet")

	}

	// Number fix
}

func properNounSpans(d *lingo.Dependency) (retVal []span) {
	start := -1
	end := -1
	for i, a := range d.AnnotatedSentence {
		if lingo.IsProperNoun(a.POSTag) {
			if start == -1 {
				start = i
				end = i + 1
			} else {
				end = i + 1
			}
		} else {
			if end == -1 {
				end = i
			}

			if start > -1 {
				s := makeSpan(start, end)
				retVal = append(retVal, s)
			}

			start = -1
			end = -1
		}
	}

	if start > -1 {
		s := makeSpan(start, len(d.AnnotatedSentence))
		retVal = append(retVal, s)
	}
	return
}
