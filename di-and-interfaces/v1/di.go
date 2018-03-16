package main

import (
"fmt"
"os"
"io"
)

func Greet(writer io.Writer, name string) {
	fmt.Fprintf(writer,"Hello, %s", name)
}

func main() {
	Greet(os.Stdout, "Elodie")
}