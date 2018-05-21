package main

import (
	"log"
	"net/http"
)

// InMemoryPlayerStore collects data about players in memory
type InMemoryPlayerStore struct{}

// RecordWin will record a player's win
func (i *InMemoryPlayerStore) RecordWin(name string) {
}

// GetPlayerScore retrieves scores for a given player
func (i *InMemoryPlayerStore) GetPlayerScore(name string) int {
	return 123
}

func main() {
	server := &PlayerServer{&InMemoryPlayerStore{}}

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
