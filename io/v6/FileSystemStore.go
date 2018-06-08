package main

import (
	"encoding/json"
	"io"
)

// FileSystemPlayerStore stores players in the filesystem
type FileSystemPlayerStore struct {
	database io.ReadWriteSeeker
}

// GetLeague returns the scores of all the players
func (f *FileSystemPlayerStore) GetLeague() (League, error) {
	f.database.Seek(0, 0)
	return NewLeague(f.database)
}

// GetPlayerScore retrieves a player's score
func (f *FileSystemPlayerStore) GetPlayerScore(name string) int {

	league, _ := f.GetLeague()
	player := league.Find(name)

	if player != nil {
		return player.Wins
	}

	return 0
}

// RecordWin will store a win for a player, incrementing wins if already known
func (f *FileSystemPlayerStore) RecordWin(name string) {
	league, _ := f.GetLeague()
	player := league.Find(name)

	if player != nil {
		player.Wins++
	} else {
		league = append(league, Player{name, 1})
	}

	f.database.Seek(0, 0)
	json.NewEncoder(f.database).Encode(league)
}
