package dep

type span struct {
	start, end int
}

func makeSpan(start, end int) span {
	if end <= start {
		panic("Impossible span created")
	}
	return span{start, end}
}

func (s span) combine(other span) span {
	start := minInt(s.start, other.start)
	end := maxInt(s.end, other.end)
	return span{start, end}
}
