package main

import (
	"bytes"
	"testing"
)

func TestCountdown(t *testing.T) {
	buffer := &bytes.Buffer{}

	Countdown(buffer)

	got := buffer.String()
	want := "5"

	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}
