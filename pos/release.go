// +build !debug

package pos

const BUILD_DEBUG = "POS TAGGER: Release Build"

var TABCOUNT uint32 = 0
var tracking = false

func tabcount() int                                   { return 0 }
func enterLoggingContext()                            {}
func leaveLoggingContext()                            {}
func logf(format string, others ...interface{})       {}
func recoverFrom(format string, attrs ...interface{}) {}

func (p *Tagger) ShowWeights() {}
func printShortcuts(p *Tagger) {}
