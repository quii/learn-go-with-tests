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
func (f *FileSystemPlayerStore) GetPlayerScore(name string) (int, error) {

	league, err := f.GetLeague()

	if err != nil{
		return 0, err
	}

	player := league.Find(name)

	if player != nil {
		return player.Wins, nil
	}

	return 0, nil
}

// RecordWin will store a win for a player, incrementing wins if already known
func (f *FileSystemPlayerStore) RecordWin(name string) error {
	league, err := f.GetLeague()

	if err != nil {
		return err
	}

	player := league.Find(name)

	if player != nil {
		player.Wins++
	} else {
		league = append(league, Player{name, 1})
	}

	f.database.Seek(0, 0)
	json.NewEncoder(f.database).Encode(league)

	return nil
}
