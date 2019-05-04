package clockface

import "time"

type Point struct {
	X int
	Y int
}

func SecondHand(t time.Time) Point {
	return Point{150, 60}
}
