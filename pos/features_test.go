// +build stanfordtags

package pos

import (
	"testing"

	"github.com/chewxy/lingo"
	"github.com/stretchr/testify/assert"
)

func TestGetFeatures(t *testing.T) {
	assert := assert.New(t)

	// test two word sentence
	s2 := lingo.AnnotatedSentence{
		lingo.AnnotationFromLexTag(lingo.Lexeme{"most", lingo.Word, -1, -1}, lingo.RBS, dummyFix{}),
		lingo.AnnotationFromLexTag(lingo.Lexeme{"populous", lingo.Word, -1, -1}, lingo.X, dummyFix{}),
	}

	featMap := getFeatures(s2, 0)
	expectedFM := featureMap{
		singleFeature{bias, ""}:                       1,
		singleFeature{ithWord_, "most"}:               1,
		tupleFeature{prevLemma_prevPOSTag, "", "X"}:   1,
		tupleFeature{prev2Lemma_prev2POSTag, "", "X"}: 1,
		singleFeature{nextWord_, "populous"}:          1,
		singleFeature{next2Word_, ""}:                 1,

		singleFeature{ithSuffix3_, "ost"}: 1,
		singleFeature{ithPrefix1_, "m"}:   1,

		singleFeature{prevPOSTag_, "X"}:                 1,
		singleFeature{prev2POSTag_, "X"}:                1,
		tupleFeature{prevPOSTag_prev2_POSTag, "X", "X"}: 1,
		tupleFeature{prevPOSTag_ithWord, "X", "most"}:   1,
		singleFeature{prevSuffix3_, ""}:                 1,
		singleFeature{nextSuffix3_, "ous"}:              1,

		singleFeature{ithShape_, "xxxx"}:  1,
		singleFeature{ithCluster_, "0"}:   1,
		singleFeature{nextCluster_, "0"}:  1,
		singleFeature{next2Cluster_, "0"}: 1,
		singleFeature{prevCluster_, "0"}:  1,
		singleFeature{prev2Cluster_, "0"}: 1,

		singleFeature{ithFlags_, "01000000010110"}:   1,
		singleFeature{nextFlags_, "00000000010110"}:  1,
		singleFeature{next2Flags_, "00000000000000"}: 1,
		singleFeature{prevFlags_, "00000000000000"}:  1,
		singleFeature{prev2Flags_, "00000000000000"}: 1,
	}
	assert.EqualValues(expectedFM, featMap, "Want: \n%v\n\nGot: \n%v", expectedFM, featMap)

	// test five word sentence
	s5 := lingo.AnnotatedSentence{
		lingo.AnnotationFromLexTag(lingo.Lexeme{"most", lingo.Word, -1, -1}, lingo.RBS, dummyFix{}),
		lingo.AnnotationFromLexTag(lingo.Lexeme{"populous", lingo.Word, -1, -1}, lingo.X, dummyFix{}),
		lingo.AnnotationFromLexTag(lingo.Lexeme{"state", lingo.Word, -1, -1}, lingo.X, dummyFix{}),
		lingo.AnnotationFromLexTag(lingo.Lexeme{"in", lingo.Word, -1, -1}, lingo.X, dummyFix{}),
		lingo.AnnotationFromLexTag(lingo.Lexeme{"America", lingo.Word, -1, -1}, lingo.X, dummyFix{}),
	}

	featMap = getFeatures(s5, 0) // no prev

	expectedFM = featureMap{
		singleFeature{bias, ""}:                       1,
		singleFeature{ithWord_, "most"}:               1,
		tupleFeature{prevLemma_prevPOSTag, "", "X"}:   1,
		tupleFeature{prev2Lemma_prev2POSTag, "", "X"}: 1,
		singleFeature{nextWord_, "populous"}:          1,
		singleFeature{next2Word_, "state"}:            1,

		singleFeature{ithSuffix3_, "ost"}: 1,
		singleFeature{ithPrefix1_, "m"}:   1,

		singleFeature{prevPOSTag_, "X"}:                 1,
		singleFeature{prev2POSTag_, "X"}:                1,
		tupleFeature{prevPOSTag_prev2_POSTag, "X", "X"}: 1,
		tupleFeature{prevPOSTag_ithWord, "X", "most"}:   1,
		singleFeature{prevSuffix3_, ""}:                 1,
		singleFeature{nextSuffix3_, "ous"}:              1,

		singleFeature{ithShape_, "xxxx"}:  1,
		singleFeature{ithCluster_, "0"}:   1,
		singleFeature{nextCluster_, "0"}:  1,
		singleFeature{next2Cluster_, "0"}: 1,
		singleFeature{prevCluster_, "0"}:  1,
		singleFeature{prev2Cluster_, "0"}: 1,

		singleFeature{ithFlags_, "01000000010110"}:   1,
		singleFeature{nextFlags_, "00000000010110"}:  1,
		singleFeature{next2Flags_, "00000000010110"}: 1,
		singleFeature{prevFlags_, "00000000000000"}:  1,
		singleFeature{prev2Flags_, "00000000000000"}: 1,
	}
	assert.EqualValues(expectedFM, featMap, "Want: \n%v\n\nGot: \n%v", expectedFM, featMap)

	featMap = getFeatures(s5, 2) // has all the feats
	expectedFM = featureMap{
		singleFeature{bias, ""}:                         1,
		singleFeature{ithWord_, "state"}:                1,
		tupleFeature{prev2Lemma_prev2POSTag, "", "RBS"}: 1,
		tupleFeature{prevLemma_prevPOSTag, "", "X"}:     1,
		singleFeature{nextWord_, "in"}:                  1,
		singleFeature{next2Word_, "america"}:            1,

		singleFeature{ithSuffix3_, "ate"}: 1,
		singleFeature{ithPrefix1_, "s"}:   1,

		singleFeature{prevPOSTag_, "X"}:                   1,
		singleFeature{prev2POSTag_, "RBS"}:                1,
		tupleFeature{prevPOSTag_prev2_POSTag, "X", "RBS"}: 1,
		tupleFeature{prevPOSTag_ithWord, "X", "state"}:    1,
		singleFeature{prevSuffix3_, "ous"}:                1,
		singleFeature{nextSuffix3_, ""}:                   1,

		singleFeature{ithShape_, "xxxx"}:  1,
		singleFeature{ithCluster_, "0"}:   1,
		singleFeature{nextCluster_, "0"}:  1,
		singleFeature{next2Cluster_, "0"}: 1,
		singleFeature{prevCluster_, "0"}:  1,
		singleFeature{prev2Cluster_, "0"}: 1,

		singleFeature{ithFlags_, "00000000010110"}:   1,
		singleFeature{nextFlags_, "01000000010110"}:  1,
		singleFeature{next2Flags_, "00000010000110"}: 1,
		singleFeature{prevFlags_, "00000000010110"}:  1,
		singleFeature{prev2Flags_, "01000000010110"}: 1,
	}
	assert.EqualValues(expectedFM, featMap, "Want: \n%v\n\nGot: \n%v", expectedFM, featMap)

	featMap = getFeatures(s5, 4) // no nexts

	expectedFM = featureMap{
		singleFeature{bias, ""}:                       1,
		singleFeature{ithWord_, "america"}:            1,
		tupleFeature{prev2Lemma_prev2POSTag, "", "X"}: 1,
		tupleFeature{prevLemma_prevPOSTag, "", "X"}:   1,
		singleFeature{nextWord_, ""}:                  1,
		singleFeature{next2Word_, ""}:                 1,

		singleFeature{ithSuffix3_, "ica"}: 1,
		singleFeature{ithPrefix1_, "A"}:   1,

		singleFeature{prevPOSTag_, "X"}:                  1,
		singleFeature{prev2POSTag_, "X"}:                 1,
		tupleFeature{prevPOSTag_prev2_POSTag, "X", "X"}:  1,
		tupleFeature{prevPOSTag_ithWord, "X", "america"}: 1,
		singleFeature{prevSuffix3_, ""}:                  1,
		singleFeature{nextSuffix3_, ""}:                  1,

		singleFeature{ithShape_, "Xxxxx"}: 1,
		singleFeature{ithCluster_, "0"}:   1,
		singleFeature{nextCluster_, "0"}:  1,
		singleFeature{next2Cluster_, "0"}: 1,
		singleFeature{prevCluster_, "0"}:  1,
		singleFeature{prev2Cluster_, "0"}: 1,

		singleFeature{ithFlags_, "00000010000110"}:   1,
		singleFeature{nextFlags_, "00000000000000"}:  1,
		singleFeature{next2Flags_, "00000000000000"}: 1,
		singleFeature{prevFlags_, "01000000010110"}:  1,
		singleFeature{prev2Flags_, "00000000010110"}: 1,
	}

	assert.EqualValues(expectedFM, featMap, "Want: \n%v\n\nGot: \n%v", expectedFM, featMap)
}
