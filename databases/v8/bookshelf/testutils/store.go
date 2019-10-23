package testutils

import (
	"errors"
	"strings"
	"time"

	"github.com/djangulo/learn-go-with-tests/databases/v8/bookshelf"
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

// List find a book by title and author. Case insensitive.
func (s *SpyStore) List(books *[]*bookshelf.Book, query string) error {
	if query == "" {
		for _, b := range s.Books {
			*books = append(*books, b)
		}
		return nil
	}
	for _, b := range s.Books {
		if strings.Contains(strings.ToLower(b.Title), query) ||
			strings.Contains(strings.ToLower(b.Author), query) {
			*books = append(*books, b)
		}
	}
	return nil
}

// Update updates a book in the SpyStore.
func (s *SpyStore) Update(book *bookshelf.Book, id int64, fields map[string]interface{}) error {
	for _, b := range s.Books {
		if (*b).ID == id {
			*book = *b
		}
	}
	if title, ok := fields["title"]; ok {
		title := title.(string)
		book.Title = title
	}
	if author, ok := fields["author"]; ok {
		author := author.(string)
		book.Author = author
	}
	return nil
}

// Delete removes a book from the SpyStore.
func (s *SpyStore) Delete(book *bookshelf.Book, id int64) error {
	err := (*s).ByID(book, id)
	if err != nil {
		return err
	}
	var idx int
	for i, b := range s.Books {
		if (*b).ID == id {
			idx = i
			break
		}
	}
	book.ID = 0
	s.Books = append(s.Books[:idx], s.Books[idx+1:]...)
	return nil
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
