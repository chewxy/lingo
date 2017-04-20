package lingo

import (
	"bytes"
	"fmt"
)

/* TAG SET */

// TagSet is a set of all the POSTags
type TagSet [MAXTAG]bool

func (ts TagSet) String() string {
	var buf bytes.Buffer
	for t, v := range ts {
		buf.WriteString(fmt.Sprintf("%v: %v\n", POSTag(t), v))
	}
	return buf.String()
}

// DependencyTypeSet is a set of all the DependencyTypes
type DependencyTypeSet [MAXDEPTYPE]bool

func (dts DependencyTypeSet) String() string {
	var buf bytes.Buffer
	for t, v := range dts {
		buf.WriteString(fmt.Sprintf("%v: %v\n", DependencyType(t), v))
	}
	return buf.String()
}
