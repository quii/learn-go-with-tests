package main

import (
	"bytes"
	"encoding/json"
	"io"
	"sort"
)

// FileSystemPlayerStore stores players in the filesystem
type FileSystemPlayerStore struct {
	database io.ReadWriteSeeker
}

// NewFileSystemPlayerStore creates a store and will initialise the database if it is a new database
func NewFileSystemPlayerStore(database io.ReadWriteSeeker) (*FileSystemPlayerStore, error) {
	buf := &bytes.Buffer{}
	length, err := io.Copy(buf, database)

	if err != nil {
		return nil, err
	}

	if length == 0 {
		json.NewEncoder(database).Encode(League{})
	}

	return &FileSystemPlayerStore{
		database: database,
	}, nil
}

// GetLeague returns the scores of all the players
func (f *FileSystemPlayerStore) GetLeague() (League, error) {
	f.database.Seek(0, 0)
	league, err := NewLeague(f.database)

	sort.Slice(league, func(i, j int) bool {
		return league[i].Wins > league[j].Wins
	})

	return league, err
}

// GetPlayerScore retrieves a player's score
func (f *FileSystemPlayerStore) GetPlayerScore(name string) (int, error) {

	league, err := f.GetLeague()

	if err != nil {
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
