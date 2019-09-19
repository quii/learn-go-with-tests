package bookshelf

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"io/ioutil"
	"log"
	"os"
)

// Store manages a bookshelf
type Store struct {
	db *pgx.Conn
}

// StoreBook will store a book
func (store *Store) StoreBook(book Book) {
	_, err := store.db.Exec(context.Background(), "insert into bookshelf.books (title, author) values ($1, $2)", book.Title, book.Author)
	if err != nil {
		log.Fatal(err)
	}
}

// GetBooks fetches all books
func (store *Store) GetBooks() ([]Book, error) {
	var books []Book

	rows, err := store.db.Query(context.Background(), "select title, author from bookshelf.books")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var title string
		var author string
		err := rows.Scan(&title, &author)
		if err != nil {
			log.Fatal(err)
		}
		books = append(books, Book{
			Title:  title,
			Author: author,
		})
	}

	return books, nil
}

// Book represents a book
type Book struct {
	Title  string
	Author string
}

// NewStore creates a new store, connecting to the db and applying db migrations
func NewStore() *Store {
	url := "postgres://postgres:learn-go-with-tests@localhost/postgres?sslmode=disable"
	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connection to database: %v\n", err)
		os.Exit(1)
	}

	_, err = conn.Exec(context.Background(), "create schema if not exists bookshelf")

	if err != nil {
		log.Fatal(err)
	}

	migration, _ := ioutil.ReadFile("0001_create_bookshelf.sql")

	_, err = conn.Exec(context.Background(), string(migration))

	if err != nil {
		log.Fatal(err)
	}

	return &Store{db: conn}
}
