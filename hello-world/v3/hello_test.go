package main

import "testing"

func TestHello(t *testing.T) {
	message := Hello("Chris")
	expected := "Hello, Chris"

	if message != expected {
		t.Errorf("expected '%s' but got '%s'", expected, message)
	}
}