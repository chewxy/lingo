package lingo

func InStringSlice(s string, l []string) bool {
	for _, v := range l {
		if s == v {
			return true
		}
	}
	return false
}

type is func(rune) bool

func StringIs(s string, f is) bool {
	for _, c := range s {
		if !f(c) {
			return false
		}
	}
	return true
}

func isAscii(r rune) bool {
	if r > 255 {
		return false
	}
	return true
}

func EqStringSlice(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
