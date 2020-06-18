package main

import (
	"encoding/json"
	"io"
)

// NewLeague creates a league from JSON.
func NewLeague(rdr io.Reader) ([]Player, error) {
	var league []Player
	err := json.NewDecoder(rdr).Decode(&league)
	return league, err
}
