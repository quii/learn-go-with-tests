package main

import "fmt"

const helloPrefix = "Hello, "

// Hello returns a personalised greeting, defaulting to Hello, world if an empty name is passed
func Hello(name string) string {
	if name == "" {
		name = "World"
	}
	return helloPrefix + name
}

func main() {
	fmt.Println(Hello("world"))
}
