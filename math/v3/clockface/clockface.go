package clockface

import (
	"math"
	"time"
)

const secondAngle = math.Pi / 30

type Hands struct {
	Hour   Vector
	Minute Vector
	Second Vector
}

type Vector struct {
	X float64
	Y float64
}

func HandsAt(t time.Time) (hands Hands) {
	return
}

func secondsInRadians(t time.Time) float64 {
	return math.Pi / (30 / float64(t.Second()))
}

func minutesInRadians(t time.Time) float64 {
	seconds := secondsInRadians(t) / 60
	minutes := math.Pi / (30 / float64(t.Minute()))
	return seconds + minutes
}

func hoursInRadians(t time.Time) float64 {
	seconds := secondsInRadians(t) / (60 * 60)
	minutes := minutesInRadians(t) / 60
	hours := math.Pi / (6 / (float64(t.Hour() % 12)))
	return seconds + minutes + hours
}

func secondHandVector(t time.Time) Vector {
    angle := secondsInRadians(t)
	x := math.Sin(angle)
	y := math.Cos(angle)

	return Vector{x, y}
}
