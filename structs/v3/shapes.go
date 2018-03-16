package main

type Rectangle struct {
	width  int
	height int
}

func Perimeter(rectangle Rectangle) (perimeter int) {
	return 2 * (rectangle.width + rectangle.height)
}

func Area(rectangle Rectangle) (area int) {
	return rectangle.width * rectangle.height
}
