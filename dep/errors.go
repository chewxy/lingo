package dep

import (
	"fmt"

	"github.com/chewxy/lingo"
)

type componentUnavailable string

func (c componentUnavailable) Error() string     { return fmt.Sprintf("%v unavailable", c) }
func (c componentUnavailable) Component() string { return string(c) }

// TarpitError is an error when the arc-standard is stuck.
// It implements GoStringer, which when called will output the state as a string.
// It also implements lingo.Sentencer, so the offending sentence can easily be retrieved
type TarpitError struct{ *configuration }

func (err TarpitError) Error() string { return "Tarpit Error" }

// NonProjective error is the error that is emitted when the dependency tree is not projective (that is to say the children cross lines)
type NonProjectiveError struct{ *lingo.Dependency }

func (err NonProjectiveError) Error() string { return "Non-projective tree" }
