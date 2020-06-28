package main

import "testing"

func TestHello(t *testing.T) {

	assertCorrectMessage := func(t *testing.T, got, want string) {
		t.Helper()
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	}

	t.Run("say hello to a person", func(t *testing.T) {
		got := Hello("Sumukha")
		want := "Hello, Sumukha"

		assertCorrectMessage(t, got, want)
	})

	t.Run("saying hello with empty string", func(t *testing.T) {
		got := Hello("")
		want := "Hello, world"

		assertCorrectMessage(t, got, want)

	})

}
