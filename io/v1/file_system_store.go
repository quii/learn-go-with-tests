package main

import (
	"io"
)

// FileSystemPlayerStore stores players in the filesystem.
type FileSystemPlayerStore struct {
	database io.ReadSeeker
}

// GetLeague returns the scores of all the players.
func (f *FileSystemPlayerStore) GetLeague() []Player {
	f.database.Seek(0, io.SeekStart)
	league, _ := NewLeague(f.database)
	return league
}
