package bookshelf_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/djangulo/learn-go-with-tests/databases/v8/bookshelf"
	"github.com/djangulo/learn-go-with-tests/databases/v8/bookshelf/testutils"
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

func TestMigrateIntegration(t *testing.T) {

	store, removeStore, err := testutils.NewTestPostgreSQLStore(false)
	if err != nil {
		fmt.Fprintf(os.Stdout, "db creation failed on test TestMigrateIntegration")
		t.FailNow()
	}
	defer removeStore()

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

}
func TestCreateIntegration(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		for _, test := range []struct {
			name, title, author string
			want                *bookshelf.Book
		}{
			{"can create", "test-title", "test-author", &bookshelf.Book{10, "test-title", "test-author"}},
		} {
			t.Run(test.name, func(t *testing.T) {
				store, removeStore, err := testutils.NewTestPostgreSQLStore(true)
				if err != nil {
					fmt.Fprintf(os.Stdout, "db creation failed on test %q", test.name)
					t.FailNow()
				}
				defer removeStore()

				book, err := bookshelf.Create(store, test.title, test.author)

				testutils.AssertNoError(t, err)
				testutils.AssertBooksEqual(t, book, test.want)
			})
		}
	})

	t.Run("failure", func(t *testing.T) {
		t.Run("cannot create a duplicate title-author", func(t *testing.T) {
			store, removeStore, err := testutils.NewTestPostgreSQLStore(true)
			if err != nil {
				fmt.Fprintf(os.Stdout, "db creation failed on test for duplicate title-author")
				t.FailNow()
			}
			defer removeStore()

			_, err = bookshelf.Create(store, "test-title", "test-author")
			if err != nil {
				t.Errorf("received error on CreateBook: %v", err)
			}

			_, err = bookshelf.Create(store, "test-title", "test-author")
			if err == nil {
				t.Error("wanted an error but didn't get one")
			}
		})
	})

}
func TestByIDIntegration(t *testing.T) {
	store, removeStore, err := testutils.NewTestPostgreSQLStore(true)
	if err != nil {
		fmt.Fprintf(os.Stdout, "db creation failed on TestByIDIntegration")
		t.FailNow()
	}
	defer removeStore()

	t.Run("success", func(t *testing.T) {
		got, err := bookshelf.ByID(store, 1)
		testutils.AssertNoError(t, err)
		want := &bookshelf.Book{1, "Alice's Adventures in Wonderland", "Lewis Carroll"}
		testutils.AssertBooksEqual(t, got, want)
	})

	t.Run("failure", func(t *testing.T) {
		for _, test := range []struct {
			name string
			in   int64
			want error
		}{
			{"not found", int64(42), bookshelf.ErrBookDoesNotExist},
			{"zero value", int64(0), bookshelf.ErrZeroValueID},
		} {
			t.Run(test.name, func(t *testing.T) {
				_, err := bookshelf.ByID(store, test.in)
				testutils.AssertError(t, err, test.want)
			})
		}
	})
}

func TestByTitleAuthorIntegration(t *testing.T) {
	store, removeStore, err := testutils.NewTestPostgreSQLStore(true)
	if err != nil {
		fmt.Fprintf(os.Stdout, "db creation failed on TestByTitleAuthorIntegration")
		t.FailNow()
	}
	defer removeStore()

	t.Run("success", func(t *testing.T) {
		got, err := bookshelf.ByTitleAuthor(store, "Alice's Adventures in Wonderland", "Lewis Carroll")
		want := &bookshelf.Book{1, "Alice's Adventures in Wonderland", "Lewis Carroll"}
		testutils.AssertNoError(t, err)
		testutils.AssertBooksEqual(t, got, want)
	})

	t.Run("failure", func(t *testing.T) {
		for _, test := range []struct {
			name, title, author string
			want                error
		}{
			{"empty title", "", "Herman Melville", bookshelf.ErrEmptyTitleField},
			{"empty author", "Moby Dick", "", bookshelf.ErrEmptyAuthorField},
			{"not found", "The DaVinci Code", "Dan Brown", bookshelf.ErrBookDoesNotExist},
		} {
			t.Run(test.name, func(t *testing.T) {
				_, err := bookshelf.ByTitleAuthor(store, test.title, test.author)
				testutils.AssertError(t, err, test.want)
			})
		}
	})
}

func TestGetOrCreateIntegration(t *testing.T) {
	store, removeStore, err := testutils.NewTestPostgreSQLStore(true)
	if err != nil {
		fmt.Fprintf(os.Stdout, "db creation failed on TestGetOrCreateIntegration")
		t.FailNow()
	}
	defer removeStore()

	t.Run("success", func(t *testing.T) {
		for _, test := range []struct {
			name, title, author string
			want                *bookshelf.Book
		}{
			{"retrieves", "Alice's Adventures in Wonderland", "Lewis Carroll", &bookshelf.Book{1, "Alice's Adventures in Wonderland", "Lewis Carroll"}},
			{"inserts", "DaVinci", "Dan Brown", &bookshelf.Book{10, "DaVinci", "Dan Brown"}},
		} {
			t.Run(test.name, func(t *testing.T) {
				book, err := bookshelf.GetOrCreate(store, test.title, test.author)
				testutils.AssertNoError(t, err)
				testutils.AssertBooksEqual(t, book, test.want)
			})
		}
	})
	t.Run("failure", func(t *testing.T) {
		for _, test := range []struct {
			name, title, author string
			want                error
		}{
			{"empty title", "", "Herman Melville", bookshelf.ErrEmptyTitleField},
			{"empty author", "Moby Dick", "", bookshelf.ErrEmptyAuthorField},
		} {
			t.Run(test.name, func(t *testing.T) {
				_, err := bookshelf.GetOrCreate(store, test.title, test.author)
				testutils.AssertError(t, err, test.want)
			})
		}
	})
}

func TestListIntegration(t *testing.T) {
	store, removeStore, err := testutils.NewTestPostgreSQLStore(true)
	if err != nil {
		fmt.Fprintf(os.Stdout, "db creation failed on TestListIntegration")
		t.FailNow()
	}
	defer removeStore()

	testBooks := []*bookshelf.Book{
		{1, "Alice's Adventures in Wonderland", "Lewis Carroll"},
		{2, "The Ball and The Cross", "G.K. Chesterton"},
		{3, "The Man Who Was Thursday", "G.K. Chesterton"},
		{4, "Moby Dick", "Herman Melville"},
		{5, "Paradise Lost", "John Milton"},
		{6, "The Tragedie of Julius Caesar", "William Shakespeare"},
		{7, "The Tragedie of Hamlet", "William Shakespeare"},
		{8, "The Tragedie of Macbeth", "William Shakespeare"},
		{9, "Romeo and Juliet", "William Shakespeare"},
	}

	for _, test := range []struct {
		name, query string
		want        []*bookshelf.Book
	}{
		{"List success: all", "", testBooks},
		{"List success: by author", "shake", testBooks[6:]},
		{"List success: by title", "alice", testBooks[0:1]},
		{"List empty if no match", "this query fails", []*bookshelf.Book{}},
	} {
		t.Run(test.name, func(t *testing.T) {
			books, err := bookshelf.List(store, test.query)
			testutils.AssertNoError(t, err)

			wantSet := make(map[string]*bookshelf.Book)
			for _, b := range test.want {
				wantSet[b.Title] = b
			}

			for _, b := range books {
				if _, ok := wantSet[b.Title]; !ok {
					t.Errorf("unwanted result %v", *b)
				}
			}
		})
	}
}

func TestUpdateIntegration(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		for _, test := range []struct {
			name     string
			fields   map[string]interface{}
			id       int64
			wantBook *bookshelf.Book
		}{
			{
				"title only",
				map[string]interface{}{"title": "Romeo & Juliet"},
				9,
				&bookshelf.Book{9, "Romeo And Juliet", "William Shakespeare"},
			},
			{
				"author only",
				map[string]interface{}{"author": "W. Shakespeare"},
				8,
				&bookshelf.Book{8, "The Tragedie of Hamlet", "W. Shakespeare"},
			},
			{
				"title+author",
				map[string]interface{}{"title": "Romeo & Juliet", "author": "W. Shakespeare"},
				9,
				&bookshelf.Book{9, "Romeo & Juliet", "W. Shakespeare"},
			},
		} {
			t.Run(test.name, func(t *testing.T) {
				store, removeStore, err := testutils.NewTestPostgreSQLStore(true)
				if err != nil {
					fmt.Fprintf(os.Stdout, "db creation failed on TestUpdateIntegration")
					t.FailNow()
				}
				defer removeStore()
				_, err = bookshelf.Update(store, test.id, test.fields)
				testutils.AssertNoError(t, err)
			})
		}

	})
	t.Run("failure", func(t *testing.T) {
		store, removeStore, err := testutils.NewTestPostgreSQLStore(true)
		if err != nil {
			fmt.Fprintf(os.Stdout, "db creation failed on TestUpdateIntegration")
			t.FailNow()
		}
		defer removeStore()
		for _, test := range []struct {
			name   string
			id     int64
			fields map[string]interface{}
			want   error
		}{
			{"empty fields", 1, nil, bookshelf.ErrEmptyFields},
			{"invalid fields", 1, map[string]interface{}{"isbn": 10}, bookshelf.ErrInvalidFields},
			{"zero value ID", 0, map[string]interface{}{"author": "doesn't matter"}, bookshelf.ErrZeroValueID},
		} {
			t.Run(test.name, func(t *testing.T) {
				_, err := bookshelf.Update(store, test.id, test.fields)
				testutils.AssertError(t, err, test.want)
			})
		}
	})
}

func TestDeleteIntegration(t *testing.T) {
	store, removeStore, err := testutils.NewTestPostgreSQLStore(true)
	if err != nil {
		fmt.Fprintf(os.Stdout, "db creation failed on TestDeleteIntegration")
		t.FailNow()
	}
	defer removeStore()

	t.Run("success", func(t *testing.T) {
		for _, test := range []struct {
			name string
			id   int64
		}{
			{"deletes", 1},
		} {
			t.Run(test.name, func(t *testing.T) {
				for _, b := range testBooks {
					disposable := new(bookshelf.Book)
					store.Create(disposable, b.Title, b.Author)
				}
				_, err := bookshelf.Delete(store, test.id)
				testutils.AssertNoError(t, err)

				_, err = bookshelf.ByID(store, test.id)
				testutils.AssertError(t, err, bookshelf.ErrBookDoesNotExist)
			})
		}
	})
	t.Run("failure", func(t *testing.T) {
		for _, test := range []struct {
			name string
			id   int64
			want error
		}{
			{"does not exist", 42, bookshelf.ErrBookDoesNotExist},
			{"zero value ID", 0, bookshelf.ErrZeroValueID},
		} {
			t.Run(test.name, func(t *testing.T) {
				_, err := bookshelf.Delete(store, test.id)
				testutils.AssertError(t, err, test.want)
			})
		}
	})

}
