package pos

import "fmt"

type componentUnavailable string

func (c componentUnavailable) Error() string     { return fmt.Sprintf("%v unavailable", c) }
func (c componentUnavailable) Component() string { return string(c) }
