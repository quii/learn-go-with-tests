package main

import "errors"

type Dict map[string]string

var NotFoundError = errors.New("could not find the word you were looking for")

func (d Dict) Search(word string) (string, error) {
	def, ok := d[word]
	if !ok {
		return "", NotFoundError
	}

	return def, nil
}
