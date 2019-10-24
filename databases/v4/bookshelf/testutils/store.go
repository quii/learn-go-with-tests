package testutils

import (
	"errors"
	"strings"
)

type migration struct {
	name   string
	stmt   string
	called int
}

// ErrNoCakeSQL no cake allowed
var ErrNoCakeSQL = errors.New("cakeSQL is not allowed")

// SpyStore store to help us spy on our functions
type SpyStore struct {
	Migrations map[string]migration
}

// ApplyMigration saves the "migration" to the store
func (s *SpyStore) ApplyMigration(name, stmt string) error {
	if strings.Contains(strings.ToLower(stmt), "cake") {
		return ErrNoCakeSQL
	}
	var m migration
	if mig, ok := s.Migrations[name]; ok {
		m = mig
		m.called++
		return nil
	}
	m = migration{
		name: name,
		stmt: stmt,
	}
	m.called++
	s.Migrations[name] = m
	return nil
}

// NewSpyStore returns a new *SpyStore
func NewSpyStore() *SpyStore {
	return &SpyStore{map[string]migration{}}
}
