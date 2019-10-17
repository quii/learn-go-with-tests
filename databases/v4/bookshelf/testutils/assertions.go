package testutils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

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

	if len(store.Migrations) == 0 {
		t.Error("no migrations in store")
	}

	m, ok := store.Migrations[name]
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

	for _, m := range store.Migrations {
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
		got = append(got, store.Migrations[m].called)
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
