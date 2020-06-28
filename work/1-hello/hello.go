package main

import "fmt"

const HelloPrefix = "Hello, "

func Hello(name string) string {
	if name == "" {
		name = "world"
	}
	return HelloPrefix + name
}

func main() {
	fmt.Println(Hello("world"))
}
