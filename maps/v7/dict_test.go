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

		assertError(t, got, ErrNotFound)
	})
}

func TestAdd(t *testing.T) {
	t.Run("new word", func(t *testing.T) {
		dict := Dict{}
		word := "test"
		def := "this is just a test"

		err := dict.Add(word, def)

		assertError(t, err, nil)
		assertDef(t, dict, word, def)
	})

	t.Run("existing word", func(t *testing.T) {
		word := "test"
		def := "this is just a test"
		dict := Dict{word: def}
		err := dict.Add(word, "new test")

		assertError(t, err, ErrWordExists)
		assertDef(t, dict, word, def)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("existing word", func(t *testing.T) {
		word := "test"
		def := "this is just a test"
		newDef := "new def"
		dict := Dict{word: def}
		err := dict.Update(word, newDef)

		assertError(t, err, nil)
		assertDef(t, dict, word, newDef)
	})

	t.Run("new word", func(t *testing.T) {
		word := "test"
		def := "this is just a test"
		dict := Dict{}

		err := dict.Update(word, def)

		assertError(t, err, ErrWordDoesNotExist)
	})
}

func TestDelete(t *testing.T) {
	word := "test"
	dict := Dict{word: "test def"}

	dict.Delete(word)

	_, err := dict.Search(word)
	if err != ErrNotFound {
		t.Errorf("Expected '%s' to be deleted", word)
	}
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
