package main

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
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
	// UP directional const
	UP uint = iota
	// DOWN directional const
	DOWN
)

var (
	// Directions holds possible direction values
	Directions = [...]string{UP: "up", DOWN: "down"}
	// ErrMigrationDirEmpty empty migration directory.
	ErrMigrationDirEmpty = errors.New("empty migration directory")
	// ErrMigrationDirNoExist migration directory does not exist.
	ErrMigrationDirNoExist = errors.New("migration directory does not exist")
)

// NewStore creates a new store, returning a connection to the db, and an
// anonymous function to remove the db connection when necessary.
func NewStore() (*Store, func()) {
	// remember to change 'secret-password' for the password you set earlier
	const connStr = "postgres://bookshelf_user:secret-password@localhost:5432/bookshelf_db?sslmode=disable"
	// if you initialized postgres with docker, the connection string will look like this
	// const connStr = "postgres://bookshelf_user:secret-password@my-postgres:5432/bookshelf_db"
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

// ApplyMigration is a wrapper around sql.DB.Exec that only returns an error.
func (s *Store) ApplyMigration(name, stmt string) error {
	_, err := s.db.Exec(stmt)
	if err != nil {
		return err
	}
	return nil
}

// MigrateUp wrapper around `migrate` that hardcodes to Directions[UP].
func MigrateUp(out io.Writer, store Storer, dir string, num int) ([]string, error) {
	return migrate(out, store, dir, num, Directions[UP])
}

// MigrateDown wrapper around `migrate` that hardcodes to Directions[DOWN].
func MigrateDown(out io.Writer, store Storer, dir string, num int) ([]string, error) {
	return migrate(out, store, dir, num, Directions[DOWN])
}

// migrate runs `num` .sql files found inside `dir`, designatied by
// `direction` which can be `up` or `down`, against the `store`. `out`
// is for reporting success or failure. migrate will abort if any error
// were to be encountered
func migrate(
	out io.Writer,
	store Storer,
	dir string,
	num int,
	direction string,
) ([]string, error) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil, ErrMigrationDirNoExist
	}

	allFiles, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	files := make([]os.FileInfo, 0)
	for _, f := range allFiles {
		if strings.HasSuffix(f.Name(), direction+".sql") {
			files = append(files, f)
		}
	}

	switch direction {
	case Directions[DOWN]:
		sort.SliceStable(files, func(i, j int) bool { return files[j].Name() < files[i].Name() })
	default:
		sort.SliceStable(files, func(i, j int) bool { return files[i].Name() < files[j].Name() })
	}

	total := len(files)
	if total == 0 {
		return nil, ErrMigrationDirEmpty
	}

	migrations := make([]string, 0)
	count := 0
	for _, file := range files {
		if num != -1 && count >= num {
			break
		}
		path := filepath.Join(dir, file.Name())
		content, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to read migration file %s, %v", file.Name(), err)
			return nil, err
		}

		fmt.Fprintf(out, "applying %d/%d: %s ", count+1, total, file.Name())
		err = store.ApplyMigration(file.Name(), string(content))
		if err != nil {
			fmt.Fprintf(out, "...FAILURE: %v\n", err)
			return nil, err
		}
		fmt.Fprint(out, "...SUCCESS\n")
		migrations = append(migrations, file.Name())
		count++
	}
	return migrations, nil
}
