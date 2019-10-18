package testutils

import (
	"errors"
	"strings"
	"time"

	"github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf"
)

type migration struct {
	created time.Time
	name    string
	stmt    string
	called  int
}

// ErrNoCakeSQL no cake allowed
var ErrNoCakeSQL = errors.New("cakeSQL is not allowed")

// SpyStore store to help us spy on our functions
type SpyStore struct {
	Migrations map[string]migration
	Books      []*bookshelf.Book
}

// CreateBook creates a new book in the store
func (s *SpyStore) CreateBook(book *bookshelf.Book, title, author string) error {
	book.ID = newID(s)
	book.Title = title
	book.Author = author
	s.Books = append(s.Books, book)
	return nil
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
		name:    name,
		stmt:    stmt,
		created: time.Now(),
	}
	m.called++
	s.Migrations[name] = m
	return nil
}

// NewSpyStore returns a new *SpyStore
func NewSpyStore() *SpyStore {
	books := make([]*bookshelf.Book, 0)
	return &SpyStore{
		Migrations: map[string]migration{},
		Books:      books,
	}
}

func newID(store *SpyStore) int64 {
	if len(store.Books) == 0 {
		return 1
	}
	var last int64
	for _, b := range store.Books {
		if b.ID > last {
			last = b.ID
		}
	}
	return last + 1
}
