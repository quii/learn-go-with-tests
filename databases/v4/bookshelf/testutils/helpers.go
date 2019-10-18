package testutils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
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
