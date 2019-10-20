package bookshelf_test

import (
	"fmt"
	"testing"

	"github.com/djangulo/learn-go-with-tests/databases/v6/bookshelf"
	"github.com/djangulo/learn-go-with-tests/databases/v6/bookshelf/testutils"
)

const queryTables = `
SELECT tablename, tableowner
FROM pg_catalog.pg_tables
WHERE
	schemaname != 'pg_catalog'
	AND
	schemaname != 'information_schema';`

type pgTable struct {
	tableOwner string `sql:"tableowner"`
	tableName  string `sql:"tablename"`
}

var (
	dbconf = bookshelf.DBConf{
		User:    bookshelf.MainDBConf.User,
		Pass:    bookshelf.MainDBConf.Pass,
		Host:    bookshelf.MainDBConf.Host,
		Port:    bookshelf.MainDBConf.Port,
		DBName:  "bookshelf_test_db",
		SSLMode: bookshelf.MainDBConf.SSLMode,
	}
)

func TestDBIntegration(t *testing.T) {

	store, removeStore, err := testutils.NewTestPostgreSQLStore(&dbconf)
	if err != nil {
		panic(err)
	}
	defer removeStore()

	t.Run("Migrate", func(t *testing.T) {

		t.Run("migrate up", func(t *testing.T) {
			_, err := bookshelf.MigrateUp(dummyWriter, store, "migrations", -1)
			if err != nil {
				t.Errorf("migration up failed: %v", err)
			}

			rows, err := store.DB.Query(queryTables)
			if err != nil {
				t.Errorf("received error querying rows: %v", err)
				t.FailNow()
			}
			defer rows.Close()

			tables := make([]pgTable, 0)
			for rows.Next() {
				var table pgTable
				if err := rows.Scan(&table.tableName, &table.tableOwner); err != nil {
					t.Errorf("error scanning row: %v", err)
					continue
				}
				tables = append(tables, table)
			}
			if err := rows.Err(); err != nil {
				t.Errorf("rows error: %v", err)
			}

			set := make(map[string]bool)
			for _, table := range tables {
				set[table.tableName] = true
			}

			if _, ok := set["books"]; !ok {
				t.Error("table books not returned")
			}
		})
		t.Run("migrate down", func(t *testing.T) {
			_, err := bookshelf.MigrateDown(dummyWriter, store, "migrations", -1)
			if err != nil {
				t.Errorf("migration down failed: %v", err)
			}

			rows, err := store.DB.Query(queryTables)
			if err != nil {
				t.Errorf("received error querying rows: %v", err)
				t.FailNow()
			}
			defer rows.Close()

			got := 0
			for rows.Next() {
				var a, b string
				if err := rows.Scan(&a, &b); err != nil {
					t.Errorf("error scanning row: %v", err)
					continue
				}
				fmt.Println(a, b)
				got++
			}
			if err := rows.Err(); err != nil {
				t.Errorf("rows error: %v", err)
			}
			if got > 0 {
				t.Errorf("got %d want 0 rows", got)
			}
		})
		t.Run("idempotency", func(t *testing.T) {
			_, err := bookshelf.MigrateDown(dummyWriter, store, "migrations", -1)
			if err != nil {
				t.Errorf("first migrate down failed: %v", err)
			}

			_, err = bookshelf.MigrateUp(dummyWriter, store, "migrations", -1)
			if err != nil {
				t.Errorf("first migrate up failed: %v", err)
			}

			_, err = bookshelf.MigrateUp(dummyWriter, store, "migrations", -1)
			if err != nil {
				t.Errorf("second migrate up failed: %v", err)
			}

			_, err = bookshelf.MigrateDown(dummyWriter, store, "migrations", -1)
			if err != nil {
				t.Errorf("second migrate down failed: %v", err)
			}

			_, err = bookshelf.MigrateDown(dummyWriter, store, "migrations", -1)
			if err != nil {
				t.Errorf("third migrate down failed: %v", err)
			}
		})
	})

	t.Run("Create", func(t *testing.T) {
		t.Run("can create a book", func(t *testing.T) {
			testutils.ResetStore(store)

			book, err := bookshelf.Create(store, "test-title", "test-author")
			if err != nil {
				t.Errorf("received error on CreateBook: %v", err)
			}
			if book.ID == 0 {
				t.Error("invalid ID received")
			}
		})

		t.Run("cannot create a duplicate title-author", func(t *testing.T) {
			testutils.ResetStore(store)

			_, err := bookshelf.Create(store, "test-title", "test-author")
			if err != nil {
				t.Errorf("received error on CreateBook: %v", err)
			}

			_, err = bookshelf.Create(store, "test-title", "test-author")
			if err == nil {
				t.Error("wanted an error but didn't get one")
			}
		})
	})
	t.Run("ByID", func(t *testing.T) {
		testutils.ResetStore(store)
		for _, b := range testBooks {
			disposable := new(bookshelf.Book)
			store.Create(disposable, b.Title, b.Author)
		}

		t.Run("ByID success", func(t *testing.T) {
			_, err := bookshelf.ByID(store, 1)
			testutils.AssertNoError(t, err)
		})

		for _, test := range []struct {
			name string
			in   int64
			want error
		}{
			{"ByID not found", int64(42), bookshelf.ErrBookDoesNotExist},
			{"ByID zero value", int64(0), bookshelf.ErrZeroValueID},
		} {
			t.Run(test.name, func(t *testing.T) {
				store := testutils.NewSpyStore(testBooks)
				_, err := bookshelf.ByID(store, test.in)
				testutils.AssertError(t, err, test.want)
			})
		}
	})

	t.Run("ByTitleAuthor", func(t *testing.T) {
		testutils.ResetStore(store)
		for _, b := range testBooks {
			disposable := new(bookshelf.Book)
			store.Create(disposable, b.Title, b.Author)
		}

		t.Run("ByTitleAuthor success", func(t *testing.T) {
			store := testutils.NewSpyStore(testBooks)

			_, err := bookshelf.ByTitleAuthor(store, testBooks[0].Title, testBooks[0].Author)
			testutils.AssertNoError(t, err)
		})

		for _, test := range []struct {
			name, title, author string
			want                error
		}{
			{"ByTitleAuthor failure empty title", "", "Herman Melville", bookshelf.ErrEmptyTitleField},
			{"ByTitleAuthor failure empty author", "Moby Dick", "", bookshelf.ErrEmptyAuthorField},
			{"ByTitleAuthor failure not found", "Moby Dick", "Herman Melville", bookshelf.ErrBookDoesNotExist},
		} {
			t.Run(test.name, func(t *testing.T) {
				store := testutils.NewSpyStore(testBooks)
				_, err := bookshelf.ByTitleAuthor(store, test.title, test.author)
				testutils.AssertError(t, err, test.want)
			})
		}
	})
}
