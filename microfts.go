// microfts
package main

import (
	//	"fmt"
	"strings"
	"unicode"
)

// based on SimpleFTS https://github.com/akrylysov/simplefts

// tokenize returns a slice of tokens for the given text.
func tokenize(text string) []string {
	return strings.FieldsFunc(text, func(r rune) bool {
		// Split on any character that is not a letter or a number.
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
}

// analyze analyzes the text and returns a slice of tokens.
func analyze(text string) []string {
	tokens := tokenize(text)
	tokens = lowercaseFilter(tokens)
	tokens = stopwordFilter(tokens)
	//	tokens = stemmerFilter(tokens)
	return tokens
}

// lowercaseFilter returns a slice of tokens normalized to lower case.
func lowercaseFilter(tokens []string) []string {
	r := make([]string, len(tokens))
	for i, token := range tokens {
		r[i] = strings.ToLower(token)
	}
	return r
}

// stopwordFilter returns a slice of tokens with stop words removed.
func stopwordFilter(tokens []string) []string {
	r := make([]string, 0, len(tokens))
	for _, token := range tokens {
		if !strings.Contains(stopwords, token) {
			r = append(r, token)
		}
	}
	return r
}

// index is an inverted index. It maps tokens to document IDs.
type index map[string][]string

// add adds documents to the index.
func (idx index) add(text, id string) {
	for _, token := range analyze(text) {
		if len(token) < 3 { //skip if small len
			continue
		}
		ids := idx[token]
		if ids != nil && ids[len(ids)-1] == id { // Don't add same ID twice.
			continue
		}
		idx[token] = append(ids, id)
	}
}

// save index to db
func (idx index) save(indb *KVBase) {
	for k, v := range idx {
		indb.MyCollection.Set([]byte(k), []byte(strings.Join(v, "\n")))
	}
}
