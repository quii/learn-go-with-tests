package main

import "errors"

// Dictionary store definitions to words
type Dictionary map[string]string

// ErrNotFound means the definition could not be found for the given word
var ErrNotFound = errors.New("could not find the word you were looking for")

// Search find a word in the dictionary
func (d Dictionary) Search(word string) (string, error) {
	definition, ok := d[word]
	if !ok {
		return "", ErrNotFound
	}

	return definition, nil
}
