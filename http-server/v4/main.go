package main

import (
	"log"
	"net/http"
)

// InMemoryPlayerStore collects data about players in memory.
type InMemoryPlayerStore struct{}

// RecordWin will record a player's win.
func (i *InMemoryPlayerStore) RecordWin(name string) {
}

// GetPlayerScore retrieves scores for a given player.
func (i *InMemoryPlayerStore) GetPlayerScore(name string) int {
	return 123
}

func main() {
	server := &PlayerServer{&InMemoryPlayerStore{}}
	log.Fatal(http.ListenAndServe(":5000", server))
}
