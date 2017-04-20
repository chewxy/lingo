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
