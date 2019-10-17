package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// Storer will hold the contract for a Store.
type Storer interface {
	ApplyMigration(name, stmt string) error
}

// Store manages a bookshelf using an *sql.DB.
type Store struct {
	db *sql.DB
}

const (
	removeTimeout = 10 * time.Second
)

// NewStore creates a new store, returning a connection to the db, and an
// anonymous function to remove the db connection when necessary
func NewStore() (*Store, func()) {
	// remember to change 'secret-password' for the password you set earlier
	const connStr = "postgres://books_user:secret-password@localhost:5432/books_db"
	// if you initialized postgres with docker, the connection string will look like this
	// const connStr = "postgres://books_user:secret-password@my-postgres:5432/books_db"
	// where 'my-postgres' is the '--name' parameter passed to the docker command

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connection to database: %v\n", err)
		os.Exit(1)
	}

	// exponential backoff
	remove := func() {
		deadline := time.Now().Add(removeTimeout)
		for tries := 0; time.Now().Before(deadline); tries++ {
			err := db.Close()
			retryIn := time.Second << uint(tries)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error closing connection to database, retrying in %v: %v\n", retryIn, err)
				time.Sleep(retryIn)
				continue
			}
			return
		}
		log.Fatalf("timeout of %v exceeded", removeTimeout)
	}

	return &Store{db: db}, remove
}

// ApplyMigration is a wrapper around sql.DB.Exec that only returns an error
func (s *Store) ApplyMigration(name, stmt string) error {
	_, err := s.db.Exec(stmt)
	if err != nil {
		return err
	}
	return nil
}
