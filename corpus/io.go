package corpus

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"io"
	"strconv"
	"strings"
)

// sortutil is a utility struct meant to sort words based on IDs
type sortutil struct {
	words []string
	ids   []int
	freqs []int
}

func (s *sortutil) Len() int           { return len(s.words) }
func (s *sortutil) Less(i, j int) bool { return s.ids[i] < s.ids[j] }
func (s *sortutil) Swap(i, j int) {
	s.words[i], s.words[j] = s.words[j], s.words[i]
	s.ids[i], s.ids[j] = s.ids[j], s.ids[i]
	if len(s.freqs) > 0 {
		s.freqs[i], s.freqs[j] = s.freqs[j], s.freqs[i]
	}
}

// ToDictWithFreq returns a simple marshalable type. Conceptually it's a JSON object with the words as the keys. The values are a pair - ID and Freq.
func ToDictWithFreq(c *Corpus) map[string]struct{ ID, Freq int } {
	retVal := make(map[string]struct{ ID, Freq int })
	for i, w := range c.words {
		retVal[w] = struct{ ID, Freq int }{i, c.frequencies[i]}
	}
	return retVal
}

// ToDict returns a marshalable dict. It returns a copy of the ID mapping.
func ToDict(c *Corpus) map[string]int {
	retVal := make(map[string]int)
	for k, v := range c.ids {
		retVal[k] = v
	}
	return retVal
}

// GobEncode implements GobEncoder for *Corpus
func (c *Corpus) GobEncode() ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)

	if err := encoder.Encode(c.words); err != nil {
		return nil, err
	}

	if err := encoder.Encode(c.ids); err != nil {
		return nil, err
	}

	if err := encoder.Encode(c.frequencies); err != nil {
		return nil, err
	}

	if err := encoder.Encode(c.maxid); err != nil {
		return nil, err
	}

	if err := encoder.Encode(c.totalFreq); err != nil {
		return nil, err
	}

	if err := encoder.Encode(c.maxWordLength); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// GobDecode implements GobDecoder for *Corpus
func (c *Corpus) GobDecode(buf []byte) error {
	b := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(b)

	if err := decoder.Decode(&c.words); err != nil {
		return err
	}

	if err := decoder.Decode(&c.ids); err != nil {
		return err
	}

	if err := decoder.Decode(&c.frequencies); err != nil {
		return err
	}

	if err := decoder.Decode(&c.maxid); err != nil {
		return err
	}

	if err := decoder.Decode(&c.totalFreq); err != nil {
		return err
	}

	if err := decoder.Decode(&c.maxWordLength); err != nil {
		return err
	}

	return nil
}

// LoadOneGram loads a 1_gram.txt file, which is a tab separated file which lists the frequency counts of words. Example:
// 		the	23135851162
// 		of	13151942776
// 		and	12997637966
// 		to	12136980858
// 		a	9081174698
// 		in	8469404971
// 		for	5933321709
func (c *Corpus) LoadOneGram(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		splits := strings.Split(line, "\t")

		if len(splits) == 0 {
			break
		}

		word := splits[0] // TODO: normalize
		count, err := strconv.Atoi(splits[1])
		if err != nil {
			return err
		}

		id := c.Add(word)
		c.frequencies[id] = count
		c.totalFreq--
		c.totalFreq += count

		wc := len([]rune(word))
		if wc > c.maxWordLength {
			c.maxWordLength = wc
		}
	}
	return nil
}
