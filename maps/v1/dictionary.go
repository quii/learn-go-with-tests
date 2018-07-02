package main

type Dictionary map[string]string

func (d Dictionary) Search(word string) string {
	return d[word]
}
