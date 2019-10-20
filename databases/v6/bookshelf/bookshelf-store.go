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

// Storer will hold the contract for a PostgreSQLStore.
type Storer interface {
	ApplyMigration(name, stmt string) error
	Create(book *Book, title string, author string) error
	ByID(book *Book, id int64) error
	ByTitleAuthor(book *Book, title string, author string) error
}

// PostgreSQLStore manages a bookshelf using an *sql.DB.
type PostgreSQLStore struct {
	DB *sql.DB
}

// Book holds book objects in the store.
type Book struct {
	ID     int64  `sql:"id"`
	Title  string `sql:"title"`
	Author string `sql:"author"`
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
	// ErrEmptyTitleField empty title field
	ErrEmptyTitleField = errors.New("empty title field")
	// ErrEmptyAuthorField empty author field
	ErrEmptyAuthorField = errors.New("empty author field")
	// ErrZeroValueID zero value ID
	ErrZeroValueID = errors.New("zero value ID")
	// ErrBookDoesNotExist book does not exist
	ErrBookDoesNotExist = errors.New("book does not exist")
	// MainDBConf holds the main database configuration
	MainDBConf DBConf
)

// DBConf holds a PostgreSQL connection configuration.
type DBConf struct {
	User, Pass, Host, Port, DBName, SSLMode string
}

// String returns the DBConf connection string.
func (d *DBConf) String() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		d.User, d.Pass, d.Host, d.Port, d.DBName, d.SSLMode)
}

// getenv is an os.Getenv extension with a default value.
func getenv(key, defaultValue string) string {
	envvar := os.Getenv(key)
	if envvar == "" {
		return defaultValue
	}
	return envvar
}

func init() {
	MainDBConf.User = getenv("POSGRES_USER", "bookshelf_user")
	MainDBConf.Pass = getenv("POSTGRES_PASSWORD", "secret-password")
	MainDBConf.Host = getenv("POSTGRES_HOST", "localhost")
	MainDBConf.Port = getenv("POSTGRES_PORT", "5432")
	MainDBConf.DBName = getenv("POSTGRES_DB", "bookshelf_db")
	MainDBConf.SSLMode = getenv("POSTGRES_SSLMODE0", "disable")
}

// NewPostgreSQLStore creates a new store, returning a connection to the db, and an
// anonymous function to remove the db connection when necessary.
func NewPostgreSQLStore(conf *DBConf) (*PostgreSQLStore, func()) {

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

	return &PostgreSQLStore{DB: db}, remove
}

// ApplyMigration is a wrapper around sql.DB.Exec that only returns an error.
func (s *PostgreSQLStore) ApplyMigration(name, stmt string) error {
	_, err := s.DB.Exec(stmt)
	if err != nil {
		return err
	}
	return nil
}

// Create inserts a new book into the postgres store.
func (s *PostgreSQLStore) Create(book *Book, title string, author string) error {
	stmt := "INSERT INTO books (title, author) VALUES ($1, $2) RETURNING id, title, author;"
	row := s.DB.QueryRow(stmt, title, author)
	err := row.Scan(&book.ID, &book.Title, &book.Author)
	if err != nil {
		return err
	}
	return nil
}

// ByID gets a book from the PostgreSQLStore by id.
func (s *PostgreSQLStore) ByID(book *Book, id int64) error {
	stmt := "SELECT id, title, author FROM books WHERE id = $1 LIMIT 1;"
	row := s.DB.QueryRow(stmt, id)
	err := row.Scan(&book.ID, &book.Title, &book.Author)
	if err != nil {
		return err
	}
	return nil
}

// ByTitleAuthor gets a book from the PostgreSQLStore by title and author.
func (s *PostgreSQLStore) ByTitleAuthor(book *Book, title, author string) error {
	stmt := "SELECT id, title, author FROM books WHERE title = $1 AND author = $2 LIMIT 1;"
	row := s.DB.QueryRow(stmt, title, author)
	err := row.Scan(&book.ID, &book.Title, &book.Author)
	if err != nil {
		return err
	}
	return nil
}

// MigrateUp wrapper around `migrate` that hardcodes to UP.
func MigrateUp(out io.Writer, store Storer, dir string, num int) ([]string, error) {
	return Migrate(out, store, dir, num, UP)
}

// MigrateDown wrapper around `migrate` that hardcodes to DOWN.
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

// Create inserts a new Book into the store.
func Create(store Storer, title, author string) (*Book, error) {
	if title == "" {
		return nil, ErrEmptyTitleField
	}
	if author == "" {
		return nil, ErrEmptyAuthorField
	}
	var book Book
	err := store.Create(&book, title, author)
	if err != nil {
		return nil, err
	}
	return &book, err
}

// ByID retrieves a book from the store by id.
func ByID(store Storer, id int64) (*Book, error) {
	if id == 0 {
		return nil, ErrZeroValueID
	}
	var book Book
	err := store.ByID(&book, id)
	if err != nil {
		return nil, ErrBookDoesNotExist
	}
	return &book, nil
}

// ByTitleAuthor retrieves a book from the store by author+title.
func ByTitleAuthor(store Storer, title, author string) (*Book, error) {
	if title == "" {
		return nil, ErrEmptyTitleField
	}
	if author == "" {
		return nil, ErrEmptyAuthorField
	}
	var book Book
	err := store.ByTitleAuthor(&book, title, author)
	if err != nil {
		return nil, ErrBookDoesNotExist
	}
	return &book, nil
}

// GetOrCreate gets a book or creates it.
func GetOrCreate(store Storer, title, author string) (*Book, error) {
	book, _ := ByTitleAuthor(store, title, author)
	if book != nil {
		return book, nil
	}
	book, err := Create(store, title, author)
	if err != nil {
		return nil, err
	}
	return book, nil
}
