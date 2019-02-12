package v1

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"
)

// SpyStore allows you to simulate a store and see how its used
type SpyStore struct {
	response string
	t        *testing.T
	ctx      context.Context
}

// Fetch returns response after a short delay
func (s *SpyStore) Fetch(ctx context.Context) (string, error) {
	data := make(chan string, 1)

	go func() {
		time.Sleep(100 * time.Millisecond)
		data <- s.response
	}()

	select {
	case msg := <-data:
		return msg, nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

type SpyResponseWriter struct {
	written bool
}

func (s *SpyResponseWriter) Header() http.Header {
	s.written = true
	return nil
}

func (s *SpyResponseWriter) Write([]byte) (int, error) {
	s.written = true
	return 0, errors.New("not implemented")
}

func (s *SpyResponseWriter) WriteHeader(statusCode int) {
	s.written = true
}
