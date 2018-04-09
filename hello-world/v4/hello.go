package main

import "fmt"

const helloPrefix = "Hello, "

// Hello returns a personalised greeting
func Hello(name string) string {
	return helloPrefix + name
}

func main() {
	fmt.Println(Hello("world"))
}
