package main

import (
	"fmt"
	"io"
	"os"
)

// Greet sends a personalised greetinb to writer
func Greet(writer io.Writer, name string) {
	fmt.Fprintf(writer, "Hello, %s", name)
}

func main() {
	Greet(os.Stdout, "Elodie")
}
