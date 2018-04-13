package racer

import (
	"net/http"
	"fmt"
)

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
