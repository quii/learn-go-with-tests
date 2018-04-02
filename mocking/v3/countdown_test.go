package main

import (
	"bytes"
	"testing"
	"time"
)

func TestCountdown(t *testing.T) {
	buffer := &bytes.Buffer{}
	spySleeper := &SpySleeper{}

	Countdown(buffer, spySleeper.Sleep)

	got := buffer.String()
	want := `5
4
3
2
1
Go!`

	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}

	if len(spySleeper.Calls) != 6 {
		t.Errorf("not enough calls to sleeper, want 6 got %d", len(spySleeper.Calls))
	}
}

type SpySleeper struct {
	Calls []time.Duration
}

func (s *SpySleeper) Sleep(duration time.Duration) {
	s.Calls = append(s.Calls, duration)
}
