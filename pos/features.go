package pos

import (
	"bytes"
	"fmt"

	"github.com/chewxy/lingo"
)

type featureType byte

//go:generate stringer -type=featureType
const (
	bias featureType = iota

	ithWord_
	nextWord_
	next2Word_

	ithSuffix3_
	ithPrefix1_

	prevPOSTag_
	prev2POSTag_
	prevSuffix3_
	nextSuffix3_

	ithShape_
	ithCluster_
	nextCluster_
	next2Cluster_
	prevCluster_
	prev2Cluster_

	ithFlags_
	nextFlags_
	next2Flags_
	prevFlags_
	prev2Flags_

	prevLemma_prevPOSTag
	prevPOSTag_ithWord
	prevPOSTag_prev2POSTag
	prev2Lemma_prev2POSTag

	MAXFEATURETYPE
)

var featCtxMap = map[featureType]contextType{
	ithWord_:   ithWord,
	nextWord_:  nextWord,
	next2Word_: next2Word,

	ithSuffix3_: ithSuffix3,
	ithPrefix1_: ithPrefix1,

	prevPOSTag_:  prevPOSTag,
	prev2POSTag_: prev2POSTag,
	prevSuffix3_: prevSuffix3,
	nextSuffix3_: nextSuffix3,

	ithShape_:     ithShape,
	ithCluster_:   ithCluster,
	nextCluster_:  nextCluster,
	next2Cluster_: next2Cluster,
	prevCluster_:  prevCluster,
	prev2Cluster_: prev2Cluster,

	ithFlags_:   ithFlags,
	nextFlags_:  nextFlags,
	next2Flags_: next2Flags,
	prevFlags_:  prevFlags,
	prev2Flags_: prev2Flags,
}

type feature interface {
	FeatType() featureType
	String() string
}

type singleFeature struct {
	featureType
	value string
}

func (sf singleFeature) FeatType() featureType { return sf.featureType }
func (sf singleFeature) String() string {
	return fmt.Sprintf("singleFeature{%v, %q}", sf.featureType, sf.value)
}

type tupleFeature struct {
	featureType
	value1 string
	value2 string
}

func (tf tupleFeature) FeatType() featureType { return tf.featureType }
func (tf tupleFeature) String() string {
	return fmt.Sprintf("tupleFeature {%v, %q, %q}", tf.featureType, tf.value1, tf.value2)
}

type featureMap map[feature]float64

func (fm featureMap) String() string {
	var buf bytes.Buffer
	for f := range fm {
		fmt.Fprintf(&buf, "%s: 1,\n", f)
	}
	return buf.String()
}

func (fm *featureMap) add(f feature) { (*fm)[f]++ }

type sfFeatures [prevLemma_prevPOSTag]singleFeature
type tfFeatures [MAXFEATURETYPE - prevLemma_prevPOSTag]tupleFeature

func fillFromContext(c contextMap) (sf sfFeatures, tf tfFeatures) {
	for i := bias; i < prevLemma_prevPOSTag; i++ {
		sf[i] = singleFeature{i, c[featCtxMap[i]]}
	}

	const last = prevLemma_prevPOSTag
	tf[prevLemma_prevPOSTag-last] = tupleFeature{prevLemma_prevPOSTag, c[prevLemma], c[prevPOSTag]}
	tf[prevPOSTag_ithWord-last] = tupleFeature{prevPOSTag_ithWord, c[prevPOSTag], c[ithWord]}
	tf[prevPOSTag_prev2POSTag-last] = tupleFeature{prevPOSTag_prev2POSTag, c[prevPOSTag], c[prev2POSTag]}
	tf[prev2Lemma_prev2POSTag-last] = tupleFeature{prev2Lemma_prev2POSTag, c[prev2Lemma], c[prev2POSTag]}
	return
}

func getFeatures(s lingo.AnnotatedSentence, i int) (sfFeatures, tfFeatures) {
	length := len(s)

	// set up context defaults
	prev2 := lingo.NullAnnotation()
	prev := lingo.NullAnnotation()
	ith := s[i]
	next := lingo.NullAnnotation()
	next2 := lingo.NullAnnotation()

	if i-1 >= 0 {
		prev = s[i-1]
	}
	if i-2 >= 0 {
		prev2 = s[i-2]
	}
	if i+1 < length {
		next = s[i+1]
	}
	if i+2 < length {
		next2 = s[i+2]
	}

	c := getContext(prev2, prev, ith, next, next2)

	return fillFromContext(c)
}
