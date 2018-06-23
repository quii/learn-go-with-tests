package main

import (
	"encoding/json"
	"os"
)

// FileSystemPlayerStore stores players in the filesystem
type FileSystemPlayerStore struct {
	database *json.Encoder
	league   League
}

// NewFileSystemPlayerStore creates a FileSystemPlayerStore
func NewFileSystemPlayerStore(file *os.File) *FileSystemPlayerStore {
	file.Seek(0, 0)
	league, _ := NewLeague(file)

	return &FileSystemPlayerStore{
		database: json.NewEncoder(&tape{file}),
		league:   league,
	}
}

// GetLeague returns the scores of all the players
func (f *FileSystemPlayerStore) GetLeague() League {
	return f.league
}

// GetPlayerScore retrieves a player's score
func (f *FileSystemPlayerStore) GetPlayerScore(name string) int {

	player := f.league.Find(name)

	if player != nil {
		return player.Wins
	}

	return 0
}

// RecordWin will store a win for a player, incrementing wins if already known
func (f *FileSystemPlayerStore) RecordWin(name string) {
	player := f.league.Find(name)

	if player != nil {
		player.Wins++
	} else {
		f.league = append(f.league, Player{name, 1})
	}

	f.database.Encode(f.league)
}
