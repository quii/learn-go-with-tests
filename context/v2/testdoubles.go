package context2

import (
	"testing"
	"time"
)

// SpyStore allows you to simulate a store and see how its used.
type SpyStore struct {
	response string
	canceled bool
	t        *testing.T
}

// Fetch returns response after a short delay.
func (s *SpyStore) Fetch() string {
	time.Sleep(100 * time.Millisecond)
	return s.response
}

// Cancel will record the call.
func (s *SpyStore) Cancel() {
	s.canceled = true
}

func (s *SpyStore) assertWasCanceled() {
	s.t.Helper()
	if !s.canceled {
		s.t.Error("store was not told to cancel")
	}
}

func (s *SpyStore) assertWasNotCanceled() {
	s.t.Helper()
	if s.canceled {
		s.t.Error("store was told to cancel")
	}
}
