package pos

import (
	"strings"
	"testing"

	"github.com/chewxy/lingo"
)

var extractContextTest = []struct {
	val string
	tag lingo.POSTag

	shape string
	pref  string
	suff  string
	flag  string
	clust string
}{
	{"TEst", lingo.ROOT_TAG, "XXxx", "T", "est", "00000000000110", "1"},
	{"TEst", lingo.X, "XXxx", "T", "est", "00000000000110", "1"},
	{"NotInClust", lingo.UNKNOWN_TAG, "XxxXxXxxxx", "N", "ust", "00000000000110", "0"},
	{"", lingo.X, "", "", "", "00000101111110", "0"},
}

func TestExtractContext(t *testing.T) {

	for i, ects := range extractContextTest {
		a := lingo.StringToAnnotation(ects.val, dummyFix{})
		a.POSTag = ects.tag

		res := extractContext(a)

		if res[0] != strings.ToLower(ects.val) {
			t.Errorf("Test %d: Expected word feature to be %q. Got %q instead", i, strings.ToLower(ects.val), res[0])
		}

		if res[2] != ects.clust {
			t.Errorf("Test %d: Expected cluster to be %q. Got %q instead", i, ects.clust, res[2])
		}

		if res[3] != ects.shape {
			t.Errorf("Test %d: Expected shape to be %q. Got %q instead", i, ects.shape, res[3])
		}

		if res[4] != ects.pref {
			t.Errorf("Test %d: Expected prefix to be %q. Got %q instead", i, ects.pref, res[4])
		}

		if res[5] != ects.suff {
			t.Errorf("Test %d: Expected suffix to be %q. Got %q instead", i, ects.suff, res[5])
		}

		if res[6] != ects.tag.String() {
			t.Errorf("Test %d: Expected postag to be %q. Got %q instead", i, ects.tag, res[6])
		}

		if res[7] != ects.flag {
			t.Errorf("Test %d: Expected flag to be %q. Got %q instead", i, ects.flag, res[7])
		}
	}

}
