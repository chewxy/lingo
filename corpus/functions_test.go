package corpus

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GenerateCorpus(t *testing.T) {
	sentenceTags := mediumSentence()
	dict := GenerateCorpus(sentenceTags)

	// testing time
	assert := assert.New(t)
	expectedWords := []string{"", "-UNKNOWN-", "-ROOT-", "President", "Bush", "on", "Tuesday", "nominated", "two", "individuals", "to", "replace", "retiring", "jurists", "federal", "courts", "in", "the", "Washington", "area", "."}

	expectedIDs := make(map[string]int)
	for i, w := range expectedWords {
		expectedIDs[w] = i
	}

	assert.Equal(expectedWords, dict.words, "Corpus known words should be the same as the manually annotated expected values")
	assert.Equal(expectedIDs, dict.ids, "IDs should be the same as expected IDs")
	assert.Equal(int64(len(expectedWords)), dict.maxid)
}

func TestViterbiSplit(t *testing.T) {
	assert := assert.New(t)
	dict := GenerateCorpus(mediumSentence())

	s2 := "twoindividuals"
	words := ViterbiSplit(s2, dict)
	assert.Equal([]string{"two", "individuals"}, words)

	s2 = "FederalCourts"
	words = ViterbiSplit(s2, dict)
	assert.Equal([]string{"federal", "courts"}, words)

	s3 := "toreplaceon"
	words = ViterbiSplit(s3, dict)
	assert.Equal([]string{"to", "replace", "on"}, words)
}

func TestCosineSimilarity(t *testing.T) {
	a := strings.Split("This is a test of cosine similarity", " ")
	b := strings.Split("This is not a test of cosine similarity", " ")

	s1 := CosineSimilarity(a, a)
	s2 := CosineSimilarity(a, b)

	if !floatEquals64(s1, 1) {
		t.Error("Expected similarity to be 1 when compared with itself")
	}
	if s2 > s1 {
		t.Error("Something went wrong with the cosine similarity algorithm")
	}

	c := strings.Split("Parramatta Road", " ")
	d := strings.Split("Parramatta Rd", " ")

	s1 = CosineSimilarity(c, c)
	s2 = CosineSimilarity(c, d)

	if !floatEquals64(s1, 1) {
		t.Error("Expected similarity to be 1 when compared with itself")
	}
	if s2 > s1 {
		t.Error("Something went wrong with the cosine similarity algorithm")
	}
}

func TestDL(t *testing.T) {
	a := "This is a test of Damerau Levenshtein"
	b := "This is not a test of Damerau Levenshtein"

	s1 := DamerauLevenshtein(a, a)
	s2 := DamerauLevenshtein(a, b)
	if s1 != 0 {
		t.Error("Expected the distance to be 0 when compared against itself. Got %d", s1)
	}

	if s2 < s1 {
		t.Error("Expected DL similarity to be greater when compared against itself")
	}

	c := "Parramatta Road"
	d := "Paramatta Rd"

	s1 = DamerauLevenshtein(c, c)
	s2 = DamerauLevenshtein(c, d)

	if s1 != 0 {
		t.Error("Expected the distance to be 0 when compared against itself. Got %d", s1)
	}
	if s2 < s1 {
		t.Error("Expected DL similarity to be greater when compared against itself")
	}
}

func TestLCP(t *testing.T) {
	assert := assert.New(t)
	lcp := LongestCommonPrefix("Hello World", "Hell yeah!")
	assert.Equal("Hell", lcp)

	lcp = LongestCommonPrefix("Hello World", "Hell yeah!", "hey there")
	assert.Equal("", lcp)

	lcp = LongestCommonPrefix()
	assert.Equal("", lcp)

	lcp = LongestCommonPrefix("OneWord")
	assert.Equal("OneWord", lcp)

	lcp = LongestCommonPrefix("foo", "foobar")
	assert.Equal("foo", lcp)
}

var parseNumTests = []struct {
	s string
	v int
}{
	{"twenty nine", 29},
	{"one hundred five", 105},
	{"five hundred twenty thousand twenty one", 520021},
}

func TestParseNumber(t *testing.T) {
	for _, pnts := range parseNumTests {
		s := strings.Split(pnts.s, " ")
		ints, err := StrsToInts(s)
		if err != nil {
			t.Error(err)
			continue
		}

		v := CombineInts(ints)
		if v != pnts.v {
			t.Errorf("Expected %q to be parsed to %d. Got %d instead", pnts.s, pnts.v, v)
		}
	}
}
