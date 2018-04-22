package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGETPlayers(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	PlayerServer(res, req)

	t.Run("hello world in response body", func(t *testing.T) {
		got := res.Body.String()
		want := "Hello, world"

		if got != want {
			t.Errorf("got '%s', want '%s'", got, want)
		}
	})

}
