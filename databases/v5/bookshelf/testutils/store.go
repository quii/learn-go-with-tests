package testutils

import (
	"errors"
	"strings"

	"github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf"
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
	Books      []*bookshelf.Book
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

// Create inserts a new *bookshelf.Book into the SpyStore.
func (s *SpyStore) Create(book *bookshelf.Book, title, author string) error {
	book.ID = newID(s)
	book.Title = title
	book.Author = author
	s.Books = append(s.Books, book)
	return nil
}

// ByID find a book by ID.
func (s *SpyStore) ByID(book *bookshelf.Book, id int64) error {
	for _, b := range s.Books {
		if b.ID == id {
			*book = *b
			return nil
		}
	}
	return bookshelf.ErrBookDoesNotExist
}

// ByTitleAuthor find a book by title and author. Case insensitive.
func (s *SpyStore) ByTitleAuthor(book *bookshelf.Book, title, author string) error {
	title, author = strings.ToLower(title), strings.ToLower(author)
	for _, b := range s.Books {
		if strings.ToLower(b.Title) == title && strings.ToLower(b.Author) == author {
			*book = *b
			return nil
		}
	}
	return bookshelf.ErrBookDoesNotExist
}

// NewSpyStore returns a new *SpyStore
func NewSpyStore(books []*bookshelf.Book) *SpyStore {
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
