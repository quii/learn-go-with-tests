package racer

import (
	"fmt"
	"net/http"
)

// Racer compares the response times of a and b, returning the fastest one
func Racer(a, b string) (winner string) {
	select {
	case <-measureResponseTime(a):
		return a
	case <-measureResponseTime(b):
		return b
	}
}

func measureResponseTime(url string) chan interface{} {
	ch := make(chan interface{})
	go func() {
		fmt.Println("getting", url)
		http.Get(url)
		ch <- true
	}()
	return ch
}
