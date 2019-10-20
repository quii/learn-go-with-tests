package bookshelf_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/djangulo/learn-go-with-tests/databases/v6/bookshelf"
	"github.com/djangulo/learn-go-with-tests/databases/v6/bookshelf/testutils"
)

func TestMigrate(t *testing.T) {

	t.Run("error on nonexistent directory", func(t *testing.T) {
		store := testutils.NewSpyStore(dummyBooks)
		_, err := bookshelf.Migrate(dummyWriter, store, "i-do-not-exist", -1, bookshelf.UP)
		testutils.AssertError(t, err, bookshelf.ErrMigrationDirNoExist)
	})

	t.Run("error on empty directory", func(t *testing.T) {
		store := testutils.NewSpyStore(dummyBooks)

		tmpdir, _, cleanup := testutils.CreateTempDir(t, "test-migrations", true)
		defer cleanup()

		_, err := bookshelf.Migrate(dummyWriter, store, tmpdir, -1, bookshelf.UP)
		testutils.AssertError(t, err, bookshelf.ErrMigrationDirEmpty)
	})

	t.Run("non-empty directory attempts to migrate", func(t *testing.T) {
		store := testutils.NewSpyStore(dummyBooks)
		tmpdir, _, cleanup := testutils.CreateTempDir(t, "test-migrations", false)
		defer cleanup()

		_, err := bookshelf.Migrate(dummyWriter, store, tmpdir, -1, bookshelf.UP)
		testutils.AssertNoError(t, err)
		testutils.AssertAllStoreMigrationCalls(t, store, 1, bookshelf.UP)
	})

	t.Run("only apply migrations in one direction", func(t *testing.T) {
		store := testutils.NewSpyStore(dummyBooks)
		tmpdir, _, cleanup := testutils.CreateTempDir(t, "test-migrations", false)
		defer cleanup()

		_, err := bookshelf.Migrate(dummyWriter, store, tmpdir, -1, bookshelf.UP)
		testutils.AssertNoError(t, err)
		testutils.AssertAllStoreMigrationCalls(t, store, 1, bookshelf.UP)
		for name := range store.Migrations {
			if strings.HasSuffix(name, "down.sql") {
				t.Errorf("Wrong direction migration applied: %s", name)
			}
		}
	})

	t.Run("up migrations should be ordered ascending", func(t *testing.T) {
		store := testutils.NewSpyStore(dummyBooks)
		tmpdir, _, cleanup := testutils.CreateTempDir(t, "test-migrations", false)
		defer cleanup()

		migrations, _ := bookshelf.Migrate(dummyWriter, store, tmpdir, -1, bookshelf.UP)
		testutils.AssertOrderAscending(t, store, migrations)
	})

	t.Run("down migrations should be ordered descending", func(t *testing.T) {
		store := testutils.NewSpyStore(dummyBooks)
		tmpdir, _, cleanup := testutils.CreateTempDir(t, "test-migrations", false)
		defer cleanup()

		migrations, _ := bookshelf.Migrate(dummyWriter, store, tmpdir, -1, bookshelf.DOWN)
		testutils.AssertOrderDescending(t, store, migrations)
	})

	t.Run("runs as many migrations as the num param, up", func(t *testing.T) {
		store := testutils.NewSpyStore(dummyBooks)
		tmpdir, _, cleanup := testutils.CreateTempDir(t, "test-migrations", false)
		defer cleanup()

		migrations, _ := bookshelf.Migrate(dummyWriter, store, tmpdir, 2, bookshelf.UP)
		testutils.AssertSliceCalls(t, store, migrations, []int{1, 1, 0})
	})

	t.Run("runs as many migrations as the num param, down", func(t *testing.T) {
		store := testutils.NewSpyStore(dummyBooks)
		tmpdir, _, cleanup := testutils.CreateTempDir(t, "test-migrations", false)
		defer cleanup()

		// keep in mind the `migrations` slice is reversed
		migrations, _ := bookshelf.Migrate(dummyWriter, store, tmpdir, 2, bookshelf.DOWN)
		testutils.AssertSliceCalls(t, store, migrations, []int{1, 1, 0})
	})

	t.Run("runs all migrations if num == -1", func(t *testing.T) {
		store := testutils.NewSpyStore(dummyBooks)
		tmpdir, _, cleanup := testutils.CreateTempDir(t, "test-migrations", false)
		defer cleanup()

		migrations, _ := bookshelf.Migrate(dummyWriter, store, tmpdir, -1, bookshelf.UP)
		testutils.AssertSliceCalls(t, store, migrations, []int{1, 1, 1})
	})

	t.Run("success output is expected", func(t *testing.T) {
		store := testutils.NewSpyStore(dummyBooks)

		tmpdir, _, cleanup := testutils.CreateTempDir(t, "test-migrations", false)
		defer cleanup()

		gotBuf := &bytes.Buffer{}
		migrations, _ := bookshelf.Migrate(gotBuf, store, tmpdir, -1, bookshelf.UP)
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
		store := testutils.NewSpyStore(dummyBooks)
		tmpdir, _, cleanup := testutils.CreateTempDir(t, "test-migrations", true)
		defer cleanup()

		tmpfile, _ := ioutil.TempFile(tmpdir, "01.cake.*.up.sql")
		tmpfile.Write([]byte("cake is superior! end pie tyranny"))
		tmpfile.Close()

		gotBuf := &bytes.Buffer{}
		_, err := bookshelf.Migrate(gotBuf, store, tmpdir, -1, bookshelf.UP)
		got := gotBuf.String()

		wantBuf := &bytes.Buffer{}
		str := fmt.Sprintf(
			"applying 1/1: %s ...FAILURE: %v\n",
			filepath.Base(tmpfile.Name()),
			testutils.ErrNoCakeSQL,
		)
		wantBuf.WriteString(str)
		want := wantBuf.String()

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
		testutils.AssertError(t, err, testutils.ErrNoCakeSQL)
	})
}

var dummyWriter = &bytes.Buffer{}
