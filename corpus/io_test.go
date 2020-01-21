package corpus

import (
	"bytes"
	"encoding/gob"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCorpusGob(t *testing.T) {
	buf := new(bytes.Buffer)

	c := New()
	c.Add("Hello")
	c.Add("World")

	helloID, _ := c.Id("Hello")
	worldID, _ := c.Id("World")

	encoder := gob.NewEncoder(buf)
	decoder := gob.NewDecoder(buf)

	if err := encoder.Encode(c); err != nil {
		t.Fatal(err)
	}

	c2 := New()
	if err := decoder.Decode(c2); err != nil {
		t.Fatal(err)
	}

	if hid, ok := c2.Id("Hello"); !ok || (ok && hid != helloID) {
		t.Errorf("\"Hello\" not found after decoding.")
	}

	if wid, ok := c2.Id("World"); !ok || (ok && wid != worldID) {
		t.Errorf("\"World\" not found after decoding.")
	}
}

func TestCorpusToDict(t *testing.T) {
	assert := assert.New(t)
	c, _ := Construct(WithWords([]string{"World", "Hello", "World"}))

	d := ToDict(c)
	c2, err := Construct(FromDict(d))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(c.words, c2.words, "Expected words to be the same")
	assert.Equal(c.ids, c2.ids, "Expected IDs to be the same")
	assert.NotEqual(c.frequencies, c2.frequencies, "Expected frequencies to not be the same")
	assert.Equal(c.maxid, c2.maxid, "Expected maxID to be the same")
	assert.NotEqual(c.totalFreq, c2.totalFreq, "Expected totalFreq to be different.")
	assert.Equal(c.maxWordLength, c2.maxWordLength, "Expected maxWordLength to be the same")
}

func TestCorpusToDictWithFreq(t *testing.T) {
	assert := assert.New(t)
	c, _ := Construct(WithWords([]string{"World", "Hello", "World"}))

	d := ToDictWithFreq(c)
	c2, err := Construct(FromDictWithFreq(d))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(c, c2)
}

func TestLoadOneGram(t *testing.T) {
	assert := assert.New(t)
	r := strings.NewReader(sample1Gram)

	c := New()
	err := c.LoadOneGram(r)
	assert.Nil(err)
	assert.Equal(10, c.Size())

	id, ok := c.Id("for")
	if !ok {
		t.Errorf("Expected \"for\" to be in corpus after loading one gram file")
	}
	assert.Equal(int(c.maxid-1), id)

}
