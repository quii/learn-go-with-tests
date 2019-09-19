package bookshelf

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"io/ioutil"
	"log"
	"os"
)

type Store struct {
	db *pgx.Conn
}

func (store *Store) StoreBook(book Book) {

}

func (store *Store) GetBooks() ([]Book, error) {

	return nil, nil
}

type Book struct {
	Title  string
	Author string
}

func NewStore() *Store {
	url := "postgres://postgres:learn-go-with-tests@localhost/postgres?sslmode=disable"
	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connection to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

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
