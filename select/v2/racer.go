package racer

import (
	"net/http"
)

// Racer compares the response times of a and b, returning the fastest one
func Racer(a, b string) (winner string) {
	select {
	case <-ping(a):
		return a
	case <-ping(b):
		return b
	}
}

func ping(url string) chan interface{} {
	ch := make(chan interface{})
	go func() {
		http.Get(url)
		ch <- true
	}()
	return ch
}
