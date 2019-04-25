package clockface

import (
	"math"
	"time"
)

type Hands struct {
	Hour   Vector
	Minute Vector
	Second Vector
}

type Vector struct {
	X int
	Y int
}

func HandsAt(t time.Time) (hands Hands) {
	return
}

func secondsInRadians(t time.Time) float64 {
	return math.Pi / (30 / (float64(t.Second())))
}

func minutesInRadians(t time.Time) float64 {
	return math.Pi / ((30 * 60) / (float64(t.Second()) +
		(60 * float64(t.Minute()))))
}

func hoursInRadians(t time.Time) float64 {
	seconds := secondsInRadians(t) / (60 * 60)
	minutes := minutesInRadians(t) / 60
	hours := math.Pi / (6 / (float64(t.Hour() % 12)))
	return seconds + minutes + hours
}


