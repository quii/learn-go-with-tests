package racer

import (
	"net/http"
	"time"
)

// Racer compares the response times of a and b, returning the fastest one.
func Racer(a, b string) (winner string) {
	aDuration := measureResponseTime(a)
	bDuration := measureResponseTime(b)

	if aDuration < bDuration {
		return a
	}

	return b
}

func measureResponseTime(url string) time.Duration {
	start := time.Now()
	resp, err := http.Get(url)
	if err == nil {
		resp.Body.Close()
	}
	return time.Since(start)
}
