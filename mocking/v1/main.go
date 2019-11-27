package main

import (
	"fmt"
	"io"
	"os"
)

// Countdown prints a countdown from 3 to out
func Countdown(out io.Writer) {
	fmt.Fprint(out, "3")
}

func main() {
	Countdown(os.Stdout)
}
