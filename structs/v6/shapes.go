package main

import "math"

type Shape interface {
	Area() float64
}

type Rectangle struct {
	width  float64
	height float64
}

func (r Rectangle) Area() (area float64) {
	return r.width * r.height
}

type Circle struct {
	radius float64
}

func (c Circle) Area() (area float64) {
	return math.Pi * c.radius * c.radius
}

func Perimeter(rectangle Rectangle) (perimeter float64) {
	return 2 * (rectangle.width + rectangle.height)
}
