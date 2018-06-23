package main

import "testing"

func TestSearch(t *testing.T) {
	dict := Dict{"test": "this is just a test"}

	got := dict.Search("test")
	want := "this is just a test"

	assertStrings(t, got, want)
}

func assertStrings(t *testing.T, got, want string) {
	t.Helper()

	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}
