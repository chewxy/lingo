package pos

import (
	"github.com/chewxy/lingo"
	"github.com/chewxy/lingo/corpus"
	"github.com/chewxy/lingo/treebank"
)

// Tagger is the object that tags an incoming channel of lexemes,
// and outputs a channel of AnnotatedSentence. Each of the Annotation
// are tagged with the POSTag
//
// The core of the Tagger is the perceptron (unexported).
//
// A large percentage of how this POS Tagger works is inspired by Mathhew Honnibal's work in SpaCy
type Tagger struct {
	*Model

	Input    chan lingo.Lexeme
	Output   chan lingo.AnnotatedSentence
	progress chan Progress

	sentences chan lingo.AnnotatedSentence

	lingo.Lemmatizer
	lingo.Stemmer
	corpus   *corpus.Corpus
	clusters map[string]lingo.Cluster // this map is safe for concurrent access because it's readonly
}

// ConsOpt is a construction option for a Tagger
type ConsOpt func(*Tagger)

// WithCorpus creates a *Tagger with an existing Corpus
func WithCorpus(c *corpus.Corpus) ConsOpt {
	fn := func(p *Tagger) {
		p.corpus = c
	}
	return fn
}

// WithLemmatizer creates a *Tagger with a lemmatizer.
// If no lemmatizer is passed into the POSTagger, then the lemmatization process will be skipped, and the POSTagger will be less accurate
func WithLemmatizer(l lingo.Lemmatizer) ConsOpt {
	fn := func(p *Tagger) {
		p.Lemmatizer = l
	}
	return fn
}

// WithStemmer creates a *Tagger with a stemmer.
// If no stemmer is passed in, then the stemming will be skipped, and the POSTagger will be less accurate
func WithStemmer(s lingo.Stemmer) ConsOpt {
	fn := func(p *Tagger) {
		p.Stemmer = s
	}
	return fn
}

// WithCluster creates a *Tagger with a brown cluster corpus (a map of strings to the brown clusters).
// If no brown cluster corpus was passed in, the cluster won't be set, and the POSTagger will be less accurate
func WithCluster(c map[string]lingo.Cluster) ConsOpt {
	fn := func(p *Tagger) {
		p.clusters = c
	}
	return fn
}

// WithModel creates a *Tagger with the specified model
func WithModel(m *Model) ConsOpt {
	fn := func(p *Tagger) {
		p.Model = m
	}
	return fn
}

// New creates a new *Tagger
func New(opts ...ConsOpt) *Tagger {
	p := &Tagger{
		Output: make(chan lingo.AnnotatedSentence),

		sentences: make(chan lingo.AnnotatedSentence),
	}

	for _, opt := range opts {
		opt(p)
	}

	if p.Model == nil {
		p.Model = &Model{perceptron: newPerceptron()}
		p.cachedTags = make(map[string]lingo.POSTag)
	}

	return p
}

// Clone() makes a copy of a POSTagger
func (p *Tagger) Clone() *Tagger {
	return &Tagger{
		Model:  p.Model,
		corpus: p.corpus,

		Output: make(chan lingo.AnnotatedSentence),

		sentences: make(chan lingo.AnnotatedSentence),

		Lemmatizer: p.Lemmatizer,
		Stemmer:    p.Stemmer,
		clusters:   p.clusters,
	}
}

// Run is used to tag a sentence. Lexemes arrive from the lexer in a channel (*Tagger.Input), and an annotated sentence is sent down the Output channel
func (p *Tagger) Run() {
	defer close(p.Output)

	go p.getSentences()

	for s := range p.sentences {
		length := len(s)
		if length == 0 {
			continue
		}
		for i, a := range s {
			tag, ok := p.shortcut(a.Lexeme)
			if !ok {
				sf, tf := getFeatures(s, i)
				tag = p.perceptron.predict(sf, tf)
			}

			p.setTag(a, tag)
		}
		p.Output <- s
	}
}

// Lemmatize implements the lingo.Lemmatize interface. It however, defers the actual doing of the job to the Lemmatizer.
func (p *Tagger) Lemmatize(a string, pt lingo.POSTag) ([]string, error) {
	if p.Lemmatizer == nil {
		return nil, componentUnavailable("lemmatizer")
	}
	return p.Lemmatizer.Lemmatize(a, pt)
}

// Stem implements the lingo.Stemmer interface. It however, defers the actual stemming to the stemmer passed in.
func (p *Tagger) Stem(a string) (string, error) {
	if p.Stemmer == nil {
		return "", componentUnavailable("stemmer")
	}
	return p.Stemmer.Stem(a)
}

// Clusters implements the lingo.AnnotationFixer interface.
func (p *Tagger) Clusters() (map[string]lingo.Cluster, error) {
	if p.clusters == nil {
		return nil, componentUnavailable("clusters")
	}
	return p.clusters, nil
}

// Progress creates and returns a channel of progress. By default the progress channel isn't created, and no progress info is sent
func (p *Tagger) Progress() <-chan Progress {
	if p.progress == nil {
		p.progress = make(chan Progress)
	}
	return p.progress
}

// Train trains a POSTagger, given a bunch of SentenceTags
func (p *Tagger) Train(sentences []treebank.SentenceTag, iterations int) {
	if p.progress != nil {
		defer func() {
			close(p.progress)
			p.progress = nil
		}()
	}

	p.fillCache(sentences)

	// Somehow sentenceTag.AnnotatedSentence() is memory leaky.
	// As a result, the more training iterations there is, the more memory is used and not released
	// hence the cache is necessary.
	cache := make(map[string]lingo.AnnotatedSentence)
	for iter := 0; iter < iterations; iter++ {
		c := 0
		n := 0
		shortcutted := 0

		var s lingo.AnnotatedSentence
		for _, sentenceTag := range sentences {
			tags := []lingo.POSTag{lingo.ROOT_TAG}
			tags = append(tags, sentenceTag.Tags...)

			var ok bool
			if s, ok = cache[sentenceTag.String()]; !ok {
				s = sentenceTag.AnnotatedSentence(p) // the fixer is used to extract cluster information, etc into the *Annotation
				cache[sentenceTag.String()] = s
			}

			length := len(s)
			if length == 0 {
				continue
			}

			for _, a := range s {
				if a == lingo.RootAnnotation() {
					continue
				}
				a.POSTag = lingo.X
			}

			for i, a := range s {
				// processing
				truth := tags[i]

				guess, ok := p.shortcut(a.Lexeme)
				if !ok {
					sf, tf := getFeatures(s, i)
					guess = p.perceptron.predict(sf, tf)
					p.perceptron.update(guess, truth, sf, tf)
				} else {
					shortcutted++
				}
				p.setTag(a, guess)

				if guess == truth {
					c++
				}
				n++
			}
		}

		if iter%150 == 0 {
			p.perceptron.average()
			logf("Averaged perceptron")
		}

		if p.progress != nil {
			p.progress <- Progress{Iter: iter, Correct: c, Count: n, ShortCutted: shortcutted}
		}

		treebank.ShuffleSentenceTag(sentences)
	}
	p.perceptron.average()
}

// LoadShortcuts allows for domain specific things to be mapped into the tagger.
func (p *Tagger) LoadShortcuts(shortcuts map[string]lingo.POSTag) {
	for shortcut, tags := range shortcuts {
		p.cachedTags[shortcut] = tags
	}
}

func (p *Tagger) fillCache(sentences []treebank.SentenceTag) {
	logf("Filling Cache with %d sentences", len(sentences))

	var counter = make(map[string]map[lingo.POSTag]int)

	for _, sentenceTag := range sentences {
		s := sentenceTag.Sentence
		tags := sentenceTag.Tags

		for i, lex := range s {
			w := lex.Value
			t := tags[i]

			_, ok := counter[w]
			if !ok {
				counter[w] = make(map[lingo.POSTag]int)
			}
			counter[w][t]++
		}
	}

	freqThresh := 30
	ambiguityThresh := 0.98

	for word, tagCounter := range counter {
		var maxTag lingo.POSTag
		var max int
		var n int
		for t, c := range tagCounter {
			if c > max {
				maxTag = t
				max = c
			}
			n += c
		}

		if n >= freqThresh && float64(max)/float64(n) >= ambiguityThresh {
			p.cachedTags[word] = maxTag
		}
	}
}

func (p *Tagger) shortcut(l lingo.Lexeme) (lingo.POSTag, bool) {
	tag, ok := lingo.POSTagShortcut(l)
	if !ok {
		tag, ok = p.cachedTags[l.Value]
	}
	return tag, ok
}

func (p *Tagger) setTag(a *lingo.Annotation, tag lingo.POSTag) {
	if a == lingo.NullAnnotation() || a == lingo.RootAnnotation() || a == lingo.StartAnnotation() {
		return
	}

	a.POSTag = tag

	if lemmas, err := p.Lemmatize(a.Value, tag); err == nil && len(lemmas) > 0 {
		// sort.Strings(lemmas)
		a.Lemma = lemmas[0]
	}

	if stem, err := p.Stem(a.Value); err == nil {
		a.Stem = stem
	}
}

// Progress is just a tuple of training progress info
type Progress struct {
	Iter, Correct, Count, ShortCutted int
}
