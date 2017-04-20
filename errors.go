package lingo

type componentUnavailable interface {
	error
	Component() string
}
