// +build !debug

package dep

const BUILD_DEBUG = "PARSER: RELEASE BUILD"
const BUILD_DIAG = "Non-Diagnostic Build"

const DEBUG = false

var READMEMSTATS = false

var TABCOUNT uint32 = 0

func enterLoggingContext() {}

func leaveLoggingContext() {}

func logTrainingProgress(iteration, correct, total, length, possibles int) {}

func logMemStats() {}

func logf(format string, others ...interface{}) {}

func recoverFrom(format string, attrs ...interface{}) {}

func (d *Parser) SprintFeatures(feature []int) string { return "" }

func SprintScores(scores []float64, ts []transition) string { return "" }
