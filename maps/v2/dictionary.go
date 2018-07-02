package main

import "errors"

type Dictionary map[string]string

var NotFoundError = errors.New("could not find the word you were looking for")

func (d Dictionary) Search(word string) (string, error) {
	def, ok := d[word]
	if !ok {
		return "", NotFoundError
	}

	return def, nil
}
