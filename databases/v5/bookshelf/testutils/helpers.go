package testutils

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf"
	_ "github.com/lib/pq" // need import
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

// NewTestStore returns a test database with the conf settings.
func NewTestStore(conf *bookshelf.DBConf) (*bookshelf.Store, func(), error) {
	main, removeMain := bookshelf.NewStore(&bookshelf.MainDBConf)

	_, err := main.DB.Exec(
		fmt.Sprintf("CREATE DATABASE %s OWNER %s;",
			conf.DBName,
			bookshelf.MainDBConf.User,
		),
	)
	if err != nil {
		return nil, nil, err
	}

	testDB, err := sql.Open("postgres", conf.String())
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
					conf.DBName,
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
			_, err := main.DB.Exec(fmt.Sprintf("DROP DATABASE %s;", conf.DBName))
			if err != nil {
				fmt.Fprintf(
					os.Stderr,
					"error dropping test database %q, retrying in %v: %v\n",
					conf.DBName,
					retryIn,
					err,
				)
				time.Sleep(retryIn)
				continue
			}
			break
		}
		removeMain()
	}
	return &bookshelf.Store{DB: testDB}, remove, nil
}

// ResetStore refreshes the store for each test.
func ResetStore(store *bookshelf.Store) error {
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
