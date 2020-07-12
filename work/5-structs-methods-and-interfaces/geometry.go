package geometry

import "math"

type Shape interface {
	Area() float64
}

type Rectangle struct {
	width  float64
	height float64
}

func (r Rectangle) Area() (area float64) {
	area = r.height * r.width
	return area
}

type Triangle struct {
	base   float64
	height float64
}

func (t Triangle) Area() (area float64) {
	area = 0.5 * t.base * t.height
	return area
}

type Circle struct {
	radius float64
}

func (c Circle) Area() (area float64) {
	area = math.Pi * c.radius * c.radius
	return area
}

func Perimeter(rectangle Rectangle) (perimeter float64) {
	perimeter = 2 * (rectangle.height + rectangle.width)
	return perimeter
}

func Area(rectangle Rectangle) (area float64) {
	area = rectangle.height * rectangle.width
	return area
}
