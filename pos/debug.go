// +build debug

package pos

import (
	"log"
	"strings"
	"sync/atomic"
)

const BUILD_DEBUG = "POS TAGGER: Debug Build"

var TABCOUNT uint32 = 0

var tracking = false

func tabcount() int {
	return int(atomic.LoadUint32(&TABCOUNT))
}

func enterLoggingContext() {
	atomic.AddUint32(&TABCOUNT, 1)
	tc := tabcount()
	log.SetPrefix(strings.Repeat("\t", tc))
}

func leaveLoggingContext() {
	tc := tabcount()
	tc--

	if tc < 0 {
		atomic.StoreUint32(&TABCOUNT, 0)
		tc = 0
	} else {
		atomic.StoreUint32(&TABCOUNT, uint32(tc))
	}
	log.SetPrefix(strings.Repeat("\t", tc))
}

func logf(format string, others ...interface{}) {
	log.Printf(format, others...)
}

func recoverFrom(format string, attrs ...interface{}) {
	if r := recover(); r != nil {
		log.Printf(format, attrs...)
		panic(r)
	}
}
