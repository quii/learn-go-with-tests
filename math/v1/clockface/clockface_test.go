package clockface_test

import (
	"testing"
	"time"

	"github.com/gypsydave5/learn-go-with-tests/math/v1/clockface"
)

func TestHandsAt(t *testing.T) {
	tm := time.Date(1337, time.January, 1, 6, 0, 0, 0, time.UTC)

	want := clockface.Hands{
		Hour:   clockface.Vector{X: 0, Y: -150},
		Minute: clockface.Vector{X: 0, Y: 150},
		Second: clockface.Vector{X: 0, Y: 150},
	}

	got := clockface.HandsAt(tm)

	if got != want {
		t.Errorf("Got %v, wanted %v", got, want)
	}
}
