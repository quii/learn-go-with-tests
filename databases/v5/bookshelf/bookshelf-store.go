package bookshelf

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

	_ "github.com/lib/pq" // unneeded namespace
)

// Storer will hold the contract for a Store.
type Storer interface {
	ApplyMigration(name, stmt string) error
	CreateBook(*Book, string, string) error
}

// Store manages a bookshelf using an *sql.DB.
type Store struct {
	DB *sql.DB
}

// Book holds book DB objects
type Book struct {
	ID     int64  `sql:"id"`
	Title  string `sql:"title"`
	Author string `sql:"author"`
}

// DBConf holds PostgreSQL connection info
type DBConf struct {
	User    string
	Pass    string
	Host    string
	Port    string
	DBName  string
	SSLMode string
}

// String returns the connection string
func (d *DBConf) String() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		d.User, d.Pass, d.Host, d.Port, d.DBName, d.SSLMode)
}

// getenv is an os.Getenv extension with a default value
func getenv(key, defaultValue string) string {
	envvar := os.Getenv(key)
	if envvar == "" {
		return defaultValue
	}
	return envvar
}

// MainDBConf holds the main database configuration
var MainDBConf DBConf

func init() {
	MainDBConf.User = getenv("POSGRES_USER", "bookshelf_user")
	MainDBConf.Pass = getenv("POSTGRES_PASSWORD", "secret-password")
	MainDBConf.Host = getenv("POSTGRES_HOST", "localhost")
	MainDBConf.Port = getenv("POSTGRES_PORT", "5432")
	MainDBConf.DBName = getenv("POSTGRES_DB", "bookshelf_db")
	MainDBConf.SSLMode = getenv("POSTGRES_SSLMODE0", "disable")
}

const (
	removeTimeout = 10 * time.Second
	// UP directional const
	UP = "up"
	// DOWN directional const
	DOWN = "down"
)

var (
	// ErrMigrationDirEmpty empty migration directory.
	ErrMigrationDirEmpty = errors.New("empty migration directory")
	// ErrMigrationDirNoExist migration directory does not exist.
	ErrMigrationDirNoExist = errors.New("migration directory does not exist")
)

// NewStore creates a new store, returning a connection to the db, and an
// anonymous function to remove the db connection when necessary.
func NewStore(conf *DBConf) (*Store, func()) {

	db, err := sql.Open("postgres", conf.String())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connection to database %q: %v\n", conf.DBName, err)
		os.Exit(1)
	}

	// exponential backoff
	remove := func() {
		deadline := time.Now().Add(removeTimeout)
		for tries := 0; time.Now().Before(deadline); tries++ {
			err := db.Close()
			retryIn := time.Second << uint(tries)
			if err != nil {
				fmt.Fprintf(
					os.Stderr,
					"error closing connection to database %q, retrying in %v: %v\n",
					conf.DBName,
					retryIn,
					err,
				)
				time.Sleep(retryIn)
				continue
			}
			return
		}
		log.Fatalf("timeout of %v exceeded", removeTimeout)
	}

	return &Store{DB: db}, remove
}

// ApplyMigration is a wrapper around sql.DB.Exec that only returns an error.
func (s *Store) ApplyMigration(name, stmt string) error {
	_, err := s.DB.Exec(stmt)
	if err != nil {
		return err
	}
	return nil
}

// CreateBook inserts a new Book into the database.
func (s *Store) CreateBook(book *Book, title, author string) error {
	stmt := "INSERT INTO books (title, author) VALUES ($1, $2) RETURNING id, title, author;"
	row := s.DB.QueryRow(stmt, title, author)
	err := row.Scan(&book.ID, &book.Title, &book.Author)
	if err != nil {
		return err
	}
	return nil
}

// MigrateUp wrapper around `migrate` that hardcodes to Directions[UP].
func MigrateUp(out io.Writer, store Storer, dir string, num int) ([]string, error) {
	return Migrate(out, store, dir, num, UP)
}

// MigrateDown wrapper around `migrate` that hardcodes to Directions[DOWN].
func MigrateDown(out io.Writer, store Storer, dir string, num int) ([]string, error) {
	return Migrate(out, store, dir, num, DOWN)
}

// Migrate runs `num` .sql files found inside `dir`, designatied by
// `direction` which can be `up` or `down`, against the `store`. `out`
// is for reporting success or failure. migrate will abort if any error
// were to be encountered
func Migrate(
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
	case DOWN:
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
