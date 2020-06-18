package main

// Rectangle has the dimensions of a rectangle.
type Rectangle struct {
	Width  float64
	Height float64
}

// Perimeter returns the perimeter of the rectangle.
func Perimeter(rectangle Rectangle) float64 {
	return 2 * (rectangle.Width + rectangle.Height)
}

// Area returns the area of the rectangle.
func Area(rectangle Rectangle) float64 {
	return rectangle.Width * rectangle.Height
}
