package corpus

import "testing"

var pluralizeTest = []struct {
	word, correct string
}{
	{"friend", "friends"},
	{"tomato", "tomatoes"},
	{"knife", "knives"},
	{"dwarf", "dwarves"},
	{"box", "boxes"},
	{"ox", "oxen"},
	{"man", "men"},
	{"equipment", "equipment"},
}

var singularizeTest = []struct {
	word, correct string
}{
	{"condolences", "condolence"},
	{"fish", "fish"},
	{"shoes", "shoe"},
	{"viri", "virus"},
	{"elves", "elf"},
}

func TestPluralize(t *testing.T) {
	for _, pts := range pluralizeTest {
		got := Pluralize(pts.word)
		if got != pts.correct {
			t.Errorf("Pluralizing %q failed. Want %q. Got %q instead", pts.word, pts.correct, got)
		}
	}
}

func TestSingularize(t *testing.T) {
	for _, pts := range singularizeTest {
		got := Singularize(pts.word)
		if got != pts.correct {
			t.Errorf("Singularizing %q failed. Want %q. Got %q instead", pts.word, pts.correct, got)
		}
	}
}
