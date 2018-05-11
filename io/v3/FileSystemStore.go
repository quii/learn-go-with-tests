package main

import (
	"io"
	"encoding/json"
)

// FileSystemPlayerStore stores players in the filesystem
type FileSystemPlayerStore struct {
	database io.ReadWriteSeeker
}

// GetLeague returns the scores of all the players
func (f *FileSystemPlayerStore) GetLeague() []Player {
	f.database.Seek(0, 0)
	league, _ := NewLeague(f.database)
	return league
}

// GetPlayerScore retrieves a player's score
func (f *FileSystemPlayerStore) GetPlayerScore(name string) int {

	var wins int

	for _, player := range f.GetLeague() {
		if player.Name == name {
			wins = player.Wins
			break
		}
	}

	return wins
}

func (f *FileSystemPlayerStore) RecordWin(name string) {
	league := f.GetLeague()

	for i, player := range league {
		if player.Name == name {
			league[i].Wins++
		}
	}

	f.database.Seek(0,0)
	json.NewEncoder(f.database).Encode(league)
}
