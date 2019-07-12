// Writes an SVG clockface of the current time to Stdout
package main

import (
	"os"
	"time"

	"github.com/gypsydave5/learn-go-with-tests/math/vFinal/clockface/svg"
)

func main() {
	t := time.Now()
	svg.Write(os.Stdout, t)
}
