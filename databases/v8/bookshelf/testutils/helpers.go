package testutils

import (
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/djangulo/learn-go-with-tests/databases/v8/bookshelf"
	_ "github.com/lib/pq" // blank import
)

// CreateTempDir creates a temporary directory. The empty flag determines
// whether it's prepopulated with sample migration files. It returns
// the directory name, the filenames in it, and a cleanup function.
// Test using it will fail if an error is raised.
func CreateTempDir(
	t *testing.T,
	name string,
	empty bool,
) (string, []string, func()) {
	t.Helper()

	tmpdir, err := ioutil.TempDir("", name)
	if err != nil {
		fmt.Println(err)
		os.RemoveAll(tmpdir)
		t.FailNow()
	}

	filenames := make([]string, 0)
	if !empty {
		for _, filename := range []string{
			"01.*.down.sql",
			"01.*.up.sql",
			"02.*.down.sql",
			"02.*.up.sql",
			"03.*.down.sql",
			"03.*.up.sql",
		} {
			tmpfile, err := ioutil.TempFile(tmpdir, filename)
			if err != nil {
				fmt.Println(err)
				t.FailNow()
			}
			filenames = append(filenames, filepath.Base(tmpfile.Name()))

			if _, err := tmpfile.Write([]byte(filename + " SQL content")); err != nil {
				tmpfile.Close()
				fmt.Println(err)
				t.FailNow()
			}
			if err := tmpfile.Close(); err != nil {
				fmt.Println(err)
				t.FailNow()
			}
		}
	}

	cleanup := func() {
		os.RemoveAll(tmpdir)
	}
	return tmpdir, filenames, cleanup

}

// TestDBRegistry keeps track of all the different test DB configs.
type TestDBRegistry struct {
	Databases map[string]*bookshelf.DBConf
	Prefix    string
}

// Add adds a new DBConfig to the registry.
func (t *TestDBRegistry) Add(conf *bookshelf.DBConf) string {
	rand.Seed(time.Now().UnixNano())
	dbname := (*t).Prefix + "_" + randString(20)

	(*conf).DBName = dbname
	(*t).Databases[dbname] = conf

	return dbname
}

// Remove drops a DBConfig from the registry
func (t *TestDBRegistry) Remove(dbname string) {
	if _, ok := (*t).Databases[dbname]; ok {
		delete((*t).Databases, dbname)
	}
}

// randString returns a random alphanumeric string of length n.
func randString(n int) string {
	b := make([]rune, n)
	for i := 0; i < n; i++ {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

var (
	chars = []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	// ActiveTestDBRegistry singleton to hold active test databases.
	ActiveTestDBRegistry = &TestDBRegistry{
		Databases: map[string]*bookshelf.DBConf{},
		Prefix:    "bookshelf_test_db",
	}
)

// NewTestPostgreSQLStore returns a test database with the bookshelf.MainDBConf settings
// and a randomly generated name.
func NewTestPostgreSQLStore(migrate bool) (*bookshelf.PostgreSQLStore, func(), error) {
	main, removeMain := bookshelf.NewPostgreSQLStore(&bookshelf.MainDBConf)
	dbconf := &bookshelf.DBConf{
		User:    bookshelf.MainDBConf.User,
		Pass:    bookshelf.MainDBConf.Pass,
		Host:    bookshelf.MainDBConf.Host,
		Port:    bookshelf.MainDBConf.Port,
		SSLMode: bookshelf.MainDBConf.SSLMode,
	}

	dbname := ActiveTestDBRegistry.Add(dbconf)

	_, err := main.DB.Exec(
		fmt.Sprintf("CREATE DATABASE %s OWNER %s;",
			dbname,
			bookshelf.MainDBConf.User,
		),
	)
	if err != nil {
		return nil, nil, err
	}

	testDB, err := sql.Open("postgres", dbconf.String())
	if err != nil {
		return nil, nil, err
	}

	remove := func() {
		closeDeadline := time.Now().Add(5 * time.Second)
		dropDeadline := time.Now().Add(10 * time.Second)
		for tries := 0; time.Now().Before(closeDeadline); tries++ {
			retryIn := time.Second << uint(tries)
			err := testDB.Close()
			if err != nil {
				fmt.Fprintf(
					os.Stderr,
					"error closing test database %q, retrying in %v: %v\n",
					dbname,
					retryIn,
					err,
				)
				time.Sleep(retryIn)
				continue
			}
			break
		}
		for tries := 0; time.Now().Before(dropDeadline); tries++ {
			retryIn := time.Second << uint(tries)
			_, err := main.DB.Exec(fmt.Sprintf("DROP DATABASE %s;", dbname))
			if err != nil {
				fmt.Fprintf(
					os.Stderr,
					"error dropping test database %q, retrying in %v: %v\n",
					dbname,
					retryIn,
					err,
				)
				time.Sleep(retryIn)
				continue
			}
			break
		}
		ActiveTestDBRegistry.Remove(dbname)
		removeMain()
	}

	store := bookshelf.PostgreSQLStore{DB: testDB}
	if migrate {
		dummyWriter := &bytes.Buffer{}
		bookshelf.MigrateUp(dummyWriter, &store, "migrations/test", -1)
	}

	return &store, remove, nil
}

// ResetStore refreshes the store for each test.
func ResetStore(store *bookshelf.PostgreSQLStore) error {
	var err error
	_, err = bookshelf.MigrateDown(dummyWriter, store, "migrations", -1)
	if err != nil {
		return err
	}

	_, err = bookshelf.MigrateUp(dummyWriter, store, "migrations", -1)
	if err != nil {
		return err
	}

	return nil
}
