package main

import (
	"testing"
)

func TestSearch(t *testing.T) {
	dict := map[string]string{"test": "this is just a test"}

	t.Run("known word", func(t *testing.T) {
		got, _ := Search(dict, "test")
		want := "this is just a test"

		assertStrings(t, got, want)
	})

	t.Run("unknown word", func(t *testing.T) {
		_, got := Search(dict, "unknown")

		assertError(t, got, NotFoundError)
	})
}

func assertStrings(t *testing.T, got, want string) {
	t.Helper()

	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}

func assertError(t *testing.T, got, want error) {
	t.Helper()

	if got != want {
		t.Errorf("got error '%s' want '%s'", got, want)
	}
}
