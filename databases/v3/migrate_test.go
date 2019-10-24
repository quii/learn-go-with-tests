package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

type migration struct {
	name   string
	stmt   string
	called int
}

type SpyStore struct {
	migrations map[string]migration
}

var errNoCakeSQL = errors.New("cakeSQL is not allowed")

// ApplyMigration saves the "migration" to the store
func (s *SpyStore) ApplyMigration(name, stmt string) error {
	if strings.Contains(strings.ToLower(stmt), "cake") {
		return errNoCakeSQL
	}
	var m migration
	if mig, ok := s.migrations[name]; ok {
		m = mig
		m.called++
		return nil
	}
	m = migration{
		name: name,
		stmt: stmt,
	}
	m.called++
	s.migrations[name] = m
	return nil
}

func NewSpyStore() *SpyStore {
	return &SpyStore{map[string]migration{}}
}

func TestMigrate(t *testing.T) {

	t.Run("error on nonexistent directory", func(t *testing.T) {
		store := NewSpyStore()
		_, err := migrate(dummyWriter, store, "i-do-not-exist", -1, UP)
		AssertError(t, err, ErrMigrationDirNoExist)
	})

	t.Run("error on empty directory", func(t *testing.T) {
		store := NewSpyStore()

		tmpdir, _, cleanup := CreateTempDir(t, "test-migrations", true)
		defer cleanup()

		_, err := migrate(dummyWriter, store, tmpdir, -1, UP)
		AssertError(t, err, ErrMigrationDirEmpty)
	})

	t.Run("non-empty directory attempts to migrate", func(t *testing.T) {
		store := NewSpyStore()
		tmpdir, _, cleanup := CreateTempDir(t, "test-migrations", false)
		defer cleanup()

		_, err := migrate(dummyWriter, store, tmpdir, -1, UP)
		AssertNoError(t, err)
		AssertAllStoreMigrationCalls(t, store, 1, UP)
	})

	t.Run("only apply migrations in one direction", func(t *testing.T) {
		store := NewSpyStore()
		tmpdir, _, cleanup := CreateTempDir(t, "test-migrations", false)
		defer cleanup()

		_, err := migrate(dummyWriter, store, tmpdir, -1, UP)
		AssertNoError(t, err)
		AssertAllStoreMigrationCalls(t, store, 1, UP)
		for name := range store.migrations {
			if strings.HasSuffix(name, "down.sql") {
				t.Errorf("Wrong direction migration applied: %s", name)
			}
		}
	})

	t.Run("up migrations should be ordered ascending", func(t *testing.T) {
		store := NewSpyStore()
		tmpdir, _, cleanup := CreateTempDir(t, "test-migrations", false)
		defer cleanup()

		migrations, _ := migrate(dummyWriter, store, tmpdir, -1, UP)
		AssertOrderAscending(t, store, migrations)
	})

	t.Run("down migrations should be ordered descending", func(t *testing.T) {
		store := NewSpyStore()
		tmpdir, _, cleanup := CreateTempDir(t, "test-migrations", false)
		defer cleanup()

		migrations, _ := migrate(dummyWriter, store, tmpdir, -1, DOWN)
		AssertOrderDescending(t, store, migrations)
	})

	t.Run("runs as many migrations as the num param, up", func(t *testing.T) {
		store := NewSpyStore()
		tmpdir, _, cleanup := CreateTempDir(t, "test-migrations", false)
		defer cleanup()

		migrations, _ := migrate(dummyWriter, store, tmpdir, 2, UP)
		AssertSliceCalls(t, store, migrations, []int{1, 1, 0})
	})

	t.Run("runs as many migrations as the num param, down", func(t *testing.T) {
		store := NewSpyStore()
		tmpdir, _, cleanup := CreateTempDir(t, "test-migrations", false)
		defer cleanup()

		// keep in mind the `migrations` slice is reversed
		migrations, _ := migrate(dummyWriter, store, tmpdir, 2, DOWN)
		AssertSliceCalls(t, store, migrations, []int{1, 1, 0})
	})

	t.Run("runs all migrations if num == -1", func(t *testing.T) {
		store := NewSpyStore()
		tmpdir, _, cleanup := CreateTempDir(t, "test-migrations", false)
		defer cleanup()

		migrations, _ := migrate(dummyWriter, store, tmpdir, -1, UP)
		AssertSliceCalls(t, store, migrations, []int{1, 1, 1})
	})

	t.Run("success output is expected", func(t *testing.T) {
		store := NewSpyStore()

		tmpdir, _, cleanup := CreateTempDir(t, "test-migrations", false)
		defer cleanup()

		gotBuf := &bytes.Buffer{}
		migrations, _ := migrate(gotBuf, store, tmpdir, -1, UP)
		got := gotBuf.String()

		total := len(migrations)

		wantBuf := &bytes.Buffer{}
		current := 1
		for _, m := range migrations {
			str := fmt.Sprintf("applying %d/%d: %s ...SUCCESS\n", current, total, m)
			wantBuf.WriteString(str)
			current++
		}
		want := wantBuf.String()

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})

	t.Run("failure output is expected", func(t *testing.T) {
		store := NewSpyStore()
		tmpdir, _, cleanup := CreateTempDir(t, "test-migrations", true)
		defer cleanup()

		tmpfile, _ := ioutil.TempFile(tmpdir, "01.cake.*.up.sql")
		tmpfile.Write([]byte("cake is superior! end pie tyranny"))
		tmpfile.Close()

		gotBuf := &bytes.Buffer{}
		_, err := migrate(gotBuf, store, tmpdir, -1, UP)
		got := gotBuf.String()

		wantBuf := &bytes.Buffer{}
		str := fmt.Sprintf(
			"applying 1/1: %s ...FAILURE: %v\n",
			filepath.Base(tmpfile.Name()),
			errNoCakeSQL,
		)
		wantBuf.WriteString(str)
		want := wantBuf.String()

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
		AssertError(t, err, errNoCakeSQL)
	})
}

// AssertError asserts error exists and is of the desired type.
func AssertError(t *testing.T, got, want error) {
	t.Helper()
	if got == nil {
		t.Error("wanted an error but didn't get one")
	}
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

// AssertNoError asserts no error is received.
func AssertNoError(t *testing.T, got error) {
	t.Helper()
	if got != nil {
		t.Errorf("got an error but didn't want one: %v", got)
	}
}

// AssertStoreMigrationCalls asserts the `name` migration exists and was
// called `num` times.
func AssertStoreMigrationCalls(t *testing.T, store *SpyStore, name string, num int) {
	t.Helper()

	if len(store.migrations) == 0 {
		t.Error("no migrations in store")
	}

	m, ok := store.migrations[name]
	if !ok {
		t.Errorf("migration %q does not exist in store", name)
	}
	if m.called != num {
		t.Errorf("got %d want %d calls migration %q", m.called, num, name)
	}
}

// AssertAllStoreMigrationCalls asserts all migrations in store for a given direction
// have a certain number of calls
func AssertAllStoreMigrationCalls(t *testing.T, store *SpyStore, num int, direction string) {
	t.Helper()

	for _, m := range store.migrations {
		if !strings.HasSuffix(m.name, direction+".sql") {
			continue
		}
		AssertStoreMigrationCalls(t, store, m.name, num)
	}
}

// AssertOrderAscending asserts the order of the `migrations` slice is ascending
// (alphabetically)
func AssertOrderAscending(t *testing.T, store *SpyStore, migrations []string) {
	t.Helper()
	for i := 0; i < len(migrations)-1; i++ {
		m0, m1 := migrations[i], migrations[i+1]
		if m0 > m1 {
			t.Errorf("wrong migration order for asc: %q before %q)", m0, m1)
		}
	}
}

// AssertOrderDescending asserts the order of the `migrations` slice is descending
// (alphabetically)
func AssertOrderDescending(t *testing.T, store *SpyStore, migrations []string) {
	t.Helper()
	for i := 0; i < len(migrations)-1; i++ {
		m0, m1 := migrations[i], migrations[i+1]
		if m0 < m1 {
			t.Errorf("wrong migration order for desc: %q before %q)", m0, m1)
		}
	}
}

// AssertSliceCalls checks the store for the `called` prorperty of migrations
func AssertSliceCalls(t *testing.T, store *SpyStore, migrations []string, want []int) {
	t.Helper()
	got := make([]int, 0)
	for _, m := range migrations {
		got = append(got, store.migrations[m].called)
	}
	for len(got) < len(want) {
		got = append(got, 0)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v calls for migrations %v", got, want, migrations)
	}
}

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

var dummyWriter = &bytes.Buffer{}
