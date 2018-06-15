package main

import "errors"

var NotFoundError = errors.New("could not find the word you were looking for")

func Search(dict map[string]string, word string) (string, error) {
	def, ok := dict[word]
	if !ok {
		return "", NotFoundError
	}

	return def, nil
}
