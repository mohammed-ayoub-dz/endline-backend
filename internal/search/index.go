package search

import (
    "strings"
    "unicode"
)

type Video struct {
    ID          int
    Title       string
    Description string
    URL         string
    Thumbnail   string
}

type Index struct {
    docs      []Video
    inverted  map[string]map[int]int 
    docLength map[int]int           
    avgDocLen float64
}

func NewIndex() *Index {
    return &Index{
        inverted:  make(map[string]map[int]int),
        docLength: make(map[int]int),
    }
}

func (idx *Index) AddDocument(doc Video) {
    idx.docs = append(idx.docs, doc)
    docID := len(idx.docs) - 1
    terms := tokenize(doc.Title + " " + doc.Description)
    termFreq := make(map[string]int)
    for _, t := range terms {
        termFreq[t]++
    }
    idx.docLength[docID] = len(terms)
    for term, freq := range termFreq {
        if _, ok := idx.inverted[term]; !ok {
            idx.inverted[term] = make(map[int]int)
        }
        idx.inverted[term][docID] = freq
    }
}

func (idx *Index) Build() {
    totalTerms := 0
    for _, length := range idx.docLength {
        totalTerms += length
    }
    if len(idx.docLength) > 0 {
        idx.avgDocLen = float64(totalTerms) / float64(len(idx.docLength))
    }
}

func tokenize(text string) []string {
    // Arabic-friendly tokenization: split on non-letter characters
    f := func(r rune) bool {
        return !unicode.IsLetter(r) && !unicode.IsNumber(r)
    }
    parts := strings.FieldsFunc(text, f)
    tokens := make([]string, 0, len(parts))
    for _, p := range parts {
        p = strings.ToLower(p)
        if len(p) > 0 {
            tokens = append(tokens, p)
        }
    }
    return tokens
}