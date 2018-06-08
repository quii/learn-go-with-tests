package main

import (
	"encoding/json"
	"io"
)

// League stores a collection of players
type League []Player

// Find tries to return a player from a league
func (l League) Find(name string) *Player {
	for i, p := range l {
		if p.Name == name {
			return &l[i]
		}
	}
	return nil
}

// NewLeague creates a league from JSON
func NewLeague(rdr io.Reader) (League, error) {
	var league []Player
	err := json.NewDecoder(rdr).Decode(&league)
	return league, err
}
