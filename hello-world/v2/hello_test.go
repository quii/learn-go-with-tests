package main

import "testing"

func TestHello(t *testing.T) {
	message := Hello()
	expected := "Hello, world"

	if message != expected {
		t.Errorf("expected '%s' but got '%s'", expected, message)
	}
}