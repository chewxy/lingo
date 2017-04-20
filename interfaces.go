package lingo

// Lemmatizer is anything that can lemmatize
type Lemmatizer interface {
	Lemmatize(string, POSTag) ([]string, error)
}

// Stemmer is anything that can stem
type Stemmer interface {
	Stem(string) (string, error)
}

// Sentencer is anything that returns an AnnotatedSentence
type Sentencer interface {
	Sentence() AnnotatedSentence
}
