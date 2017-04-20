package dep

import "github.com/chewxy/lingo"

// the features are used as columns in the matrix

// go:generate stringer type=feature -output=feature_string.go
type feature int

const (
	// first 18 are word related features
	// second 18 are POS related features
	// last 12 are label related features

	s0w feature = iota
	s1w
	s2w

	b0w
	b1w
	b2w

	s0l1w
	s0r1w
	s0l2w
	s0r2w
	s0llw
	s0rrw

	s1l1w
	s1r1w
	s1l2w
	s1r2w
	s1llw
	s1rrw

	// POS related words
	s0t
	s1t
	s2t

	b0t
	b1t
	b2t

	s0l1t
	s0r1t
	s0l2t
	s0r2t
	s0llt
	s0rrt

	s1l1t
	s1r1t
	s1l2t
	s1r2t
	s1llt
	s1rrt

	// label related
	s0l1d
	s0r1d
	s0l2d
	s0r2d
	s0lld
	s0rrd

	s1l1d
	s1r1d
	s1l2d
	s1r2d
	s1lld
	s1rrd

	MAXFEATURE
)

const (
	wordFeatsStartAt  int = int(lingo.MAXTAG) + int(lingo.MAXDEPTYPE)
	labelFeatsStartAt     = int(lingo.MAXTAG)
	posFeatsStartAt       = 0
)
