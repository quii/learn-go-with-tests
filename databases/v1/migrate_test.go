package main

import (
	"testing"
	"time"
)

type migration struct {
	created time.Time
	name    string
	stmt    string
	called  int
}

type SpyStore struct {
	migrations map[string]migration
}

func (s *SpyStore) ApplyMigration(name, stmt string) error {
	mig := s.migrations[name]
	mig.name = name
	mig.stmt = stmt
	mig.called++
	return nil
}

func NewSpyStore() *SpyStore {
	return &SpyStore{map[string]migration{}}
}

func TestMigrate(t *testing.T) {
	store := NewSpyStore()

	t.Run("error on nonexistent directory", func(t *testing.T) {
		err := migrate(store, "i-do-not-exist", -1)
		if err == nil {
			t.Error("wanted an error but didn't get one")
		}
	})

	t.Run("no error on existing directory", func(t *testing.T) {
		err := migrate(store, "migrations", -1)
		if err != nil {
			t.Errorf("got an error but didn't want one: %v", err)
		}
	})
}
