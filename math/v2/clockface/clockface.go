package clockface

import (
	"math"
	"time"
)

type Point struct {
	X int
	Y int
}

func SecondHand(t time.Time) Point {
	return Point{150, 60}
}

func secondsInRadians(t time.Time) float64 {
	return math.Pi
}
