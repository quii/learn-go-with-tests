package bookshelf_test

import (
	"fmt"
	"testing"

	"github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf"
	"github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf/testutils"
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

	store, removeStore, err := testutils.NewTestStore(&dbconf)
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

	t.Run("CreateBook", func(t *testing.T) {
		t.Run("can create a book", func(t *testing.T) {
			testutils.ResetStore(store)

			var book bookshelf.Book
			err := store.CreateBook(&book, "test-title", "test-author")
			if err != nil {
				t.Errorf("received error on CreateBook: %v", err)
			}
			if book.ID == 0 {
				t.Error("invalid ID received")
			}
		})

		t.Run("cannot create a duplicate title-author", func(t *testing.T) {
			testutils.ResetStore(store)

			var b1, b2 bookshelf.Book
			err := store.CreateBook(&b1, "test-title", "test-author")
			if err != nil {
				t.Errorf("received error on CreateBook: %v", err)
			}

			err = store.CreateBook(&b2, "test-title", "test-author")
			if err == nil {
				t.Error("wanted an error but didn't get one")
			}
		})
	})
}
