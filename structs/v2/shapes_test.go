package main

import "testing"

func TestPerimeter(t *testing.T) {
	got := Perimeter(10, 10)
	want := 40

	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}

func TestArea(t *testing.T) {
	got := Area(12, 6)
	want := 72

	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}
