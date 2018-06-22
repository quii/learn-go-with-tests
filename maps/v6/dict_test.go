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

		assertError(t, got, ErrNotFound)
	})
}

func TestAdd(t *testing.T) {
	t.Run("new word", func(t *testing.T) {
		dict := map[string]string{}
		word := "test"
		def := "this is just a test"

		err := Add(dict, word, def)

		assertError(t, err, nil)
		assertDef(t, dict, word, def)
	})

	t.Run("existing word", func(t *testing.T) {
		word := "test"
		def := "this is just a test"
		dict := map[string]string{word: def}
		err := Add(dict, word, "new test")

		assertError(t, err, ErrWordExists)
		assertDef(t, dict, word, def)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("existing word", func(t *testing.T) {
		word := "test"
		def := "this is just a test"
		newDef := "new def"
		dict := map[string]string{word: def}
		err := Update(dict, word, newDef)

		assertError(t, err, nil)
		assertDef(t, dict, word, newDef)
	})

	t.Run("new word", func(t *testing.T) {
		word := "test"
		def := "this is just a test"
		dict := map[string]string{}

		err := Update(dict, word, def)

		assertError(t, err, ErrWordDoesNotExist)
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

func assertDef(t *testing.T, dict map[string]string, word, def string) {
	t.Helper()

	got, err := Search(dict, word)
	if err != nil {
		t.Fatal("should find added word:", err)
	}

	if def != got {
		t.Errorf("got '%s' want '%s'", got, def)
	}
}
