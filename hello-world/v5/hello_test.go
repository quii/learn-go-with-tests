package main

import "testing"

func TestHello(t *testing.T) {

	assertCorrectMessage := func(expected, actual string) {
		if expected != actual {
			t.Errorf("expected '%s' but got '%s'", expected, actual)
		}
	}

	t.Run("saying hello to people", func(t *testing.T) {
		message := Hello("Chris")
		expected := "Hello, Chris"
		assertCorrectMessage(expected, message)
	})

	t.Run("say hello world when an empty string is supplied", func(t *testing.T) {
		message := Hello("")
		expected := "Hello, World"
		assertCorrectMessage(expected, message)
	})

}