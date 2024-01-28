package main

import "testing"

func TestHello(t *testing.T) {

	testcase := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "to a person",
			test: func(t *testing.T) {
				got := Hello("", "")
				want := "Hello, World"
				assertCorrectMessage(t, got, want)
			},
		}, {
			name: "empty string",
			test: func(t *testing.T) {
				got := Hello("", "")
				want := "Hello, World"
				assertCorrectMessage(t, got, want)
			},
		}, {
			name: "in Spanish",
			test: func(t *testing.T) {
				got := Hello("Elodie", spanish)
				want := "Hola, Elodie"
				assertCorrectMessage(t, got, want)
			},
		}, {
			name: "in French",
			test: func(t *testing.T) {
				got := Hello("Lauren", french)
				want := "Bonjour, Lauren"
				assertCorrectMessage(t, got, want)
			},
		},
	}

	for i := range testcase {
		tc := testcase[i]
		t.Run(tc.name, tc.test)
	}
}

func assertCorrectMessage(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}