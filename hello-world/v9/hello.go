package main

import "fmt"

const defaultName = "World"
const defaultLanguage = "english"

var prefixes = map[string]string{
	"french": "Bonjour, ",
	"spanish": "Hola, ",
	"english": "Hello, ",
}

// Hello returns a personalised greeting in a given language
func Hello(name string, language string) string {
	if len(name) == 0 {
		name = defaultName
	}

	if len(language) == 0 {
		language = defaultLanguage
	}

	return prefixes[language] + name
}

func main() {
	fmt.Println(Hello("", ""))
}
