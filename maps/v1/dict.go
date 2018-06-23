package main

type Dict map[string]string

func (d Dict) Search(word string) string {
	return d[word]
}
