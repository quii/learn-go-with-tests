package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// Storer will hold the contract for a PostgreSQLStore.
type Storer interface {
	ApplyMigration(name, stmt string) error
}

// PostgreSQLStore manages a bookshelf using an *sql.DB.
type PostgreSQLStore struct {
	db *sql.DB
}

const (
	removeTimeout = 10 * time.Second
)

// NewPostgreSQLStore creates a new store, returning a connection to the db, and an
// anonymous function to remove the db connection when necessary
func NewPostgreSQLStore() (*PostgreSQLStore, func()) {
	// remember to change 'secret-password' for the password you set earlier
	const connStr = "postgres://bookshelf_user:secret-password@localhost:5432/bookshelf_db"

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

	return &PostgreSQLStore{db: db}, remove
}

// ApplyMigration is a wrapper around sql.DB.Exec that only returns an error
func (s *PostgreSQLStore) ApplyMigration(name, stmt string) error {
	_, err := s.db.Exec(stmt)
	if err != nil {
		return err
	}
	return nil
}

func main() {}
