package main

type Rectangle struct {
	width  float64
	height float64
}

func Perimeter(rectangle Rectangle) (perimeter float64) {
	return 2 * (rectangle.width + rectangle.height)
}

func Area(rectangle Rectangle) (area float64) {
	return rectangle.width * rectangle.height
}
