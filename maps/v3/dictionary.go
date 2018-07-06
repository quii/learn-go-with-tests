package main

import "errors"

// Dictionary store definitions to words
type Dictionary map[string]string

// NotFoundError means the definition could not be found for the given word
var NotFoundError = errors.New("could not find the word you were looking for")

// Search find a word in the dictionary
func (d Dictionary) Search(word string) (string, error) {
	definition, ok := d[word]
	if !ok {
		return "", NotFoundError
	}

	return definition, nil
}

// Add inserts a word and definition into the dictionary
func (d Dictionary) Add(word, definition string) {
	d[word] = definition
}
