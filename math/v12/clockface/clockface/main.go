package main

import (
	"os"
	"time"

	"github.com/quii/learn-go-with-tests/math/v12/clockface"
)

func main() {
	t := time.Now()
	clockface.SVGWriter(os.Stdout, t)
}
