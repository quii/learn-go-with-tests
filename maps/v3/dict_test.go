package main

import (
	"testing"
)

func TestSearch(t *testing.T) {
	dict := Dict{"test": "this is just a test"}

	t.Run("known word", func(t *testing.T) {
		got, _ := dict.Search("test")
		want := "this is just a test"

		assertStrings(t, got, want)
	})

	t.Run("unknown word", func(t *testing.T) {
		_, got := dict.Search("unknown")

		assertError(t, got, NotFoundError)
	})
}

func TestAdd(t *testing.T) {
	dict := Dict{}
	word := "test"
	def := "this is just a test"

	dict.Add(word, def)

	assertDef(t, dict, word, def)
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

func assertDef(t *testing.T, dict Dict, word, def string) {
	t.Helper()

	got, err := dict.Search(word)
	if err != nil {
		t.Fatal("should find added word:", err)
	}

	if def != got {
		t.Errorf("got '%s' want '%s'", got, def)
	}
}
