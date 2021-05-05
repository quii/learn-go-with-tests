package main

import "fmt"

func Hello(name string) string {
    return "Hello, " + name
}

func main() {
    fmt.Println(Hello("world"))
}