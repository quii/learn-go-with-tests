package main

import "fmt"

const spanish = "Spanish"
const french = "French"
const helloPrefix = "Hello, "
const spanishHelloPrefix = "Hola, "
const frenchHelloPrefix = "Bonjour, "

func Hello(name string, language string) string {
	if name == "" {
		name = "World"
	}

	if language == spanish {
		return spanishHelloPrefix + name
	}

	if language == french {
		return frenchHelloPrefix + name
	}

	return helloPrefix + name
}

func main() {
	fmt.Println(Hello("world", ""))
}
