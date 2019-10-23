package bookshelf_test

import (
	"reflect"
	"testing"

	"github.com/djangulo/learn-go-with-tests/databases/v8/bookshelf"
	"github.com/djangulo/learn-go-with-tests/databases/v8/bookshelf/testutils"
)

func TestCreate(t *testing.T) {

	t.Run("creates accurately", func(t *testing.T) {
		store := testutils.NewSpyStore(dummyBooks)

		title, author := "Moby Dick", "Herman Melville"

		book, err := bookshelf.Create(store, title, author)
		testutils.AssertNoError(t, err)
		testutils.AssertBooksEqual(t, book, &bookshelf.Book{1, "Moby Dick", "Herman Melville"})
	})
	t.Run("error on empty title", func(t *testing.T) {
		store := testutils.NewSpyStore(dummyBooks)

		title, author := "", "Herman Melville"

		_, err := bookshelf.Create(store, title, author)
		testutils.AssertError(t, err, bookshelf.ErrEmptyTitleField)
	})
	t.Run("error on empty author", func(t *testing.T) {
		store := testutils.NewSpyStore(dummyBooks)

		title, author := "Moby Dick", ""

		_, err := bookshelf.Create(store, title, author)
		testutils.AssertError(t, err, bookshelf.ErrEmptyAuthorField)
	})
}

func TestByID(t *testing.T) {

	t.Run("ByID success", func(t *testing.T) {
		store := testutils.NewSpyStore(testBooks)

		book, err := bookshelf.ByID(store, testBooks[0].ID)
		testutils.AssertNoError(t, err)
		testutils.AssertBooksEqual(t, book, testBooks[0])
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

}

func TestByTitleAuthor(t *testing.T) {

	t.Run("ByTitleAuthor success", func(t *testing.T) {
		store := testutils.NewSpyStore(testBooks)

		book, err := bookshelf.ByTitleAuthor(store, testBooks[0].Title, testBooks[0].Author)
		testutils.AssertNoError(t, err)
		testutils.AssertBooksEqual(t, book, testBooks[0])
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

}

func TestGetOrCreate(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		for _, test := range []struct {
			name, title, author string
			want                *bookshelf.Book
		}{
			{"GetOrCreate retrieves from store", testBooks[1].Title, testBooks[1].Author, testBooks[1]},
			{
				"GetOrCreate insert into store",
				"Moby Dick",
				"Herman Melville",
				&bookshelf.Book{testBooks[len(testBooks)-1].ID + 1, "Moby Dick", "Herman Melville"},
			},
		} {
			t.Run(test.name, func(t *testing.T) {
				store := testutils.NewSpyStore(testBooks)
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
			{"GetOrCreate failure empty title", "", "Herman Melville", bookshelf.ErrEmptyTitleField},
			{"GetOrCreate failure empty author", "Moby Dick", "", bookshelf.ErrEmptyAuthorField},
		} {
			t.Run(test.name, func(t *testing.T) {
				store := testutils.NewSpyStore(testBooks)
				_, err := bookshelf.GetOrCreate(store, test.title, test.author)
				testutils.AssertError(t, err, test.want)
			})
		}
	})
}

func TestList(t *testing.T) {
	store := testutils.NewSpyStore(testBooks)

	for _, test := range []struct {
		name, query string
		want        []*bookshelf.Book
	}{
		{"List success: all", "", testBooks},
		{"List success: by author", "shake", testBooks[:2]},
		{"List success: by title", "old man", testBooks[2:3]},
		{"List empty if no match", "this query fails", []*bookshelf.Book{}},
	} {
		t.Run(test.name, func(t *testing.T) {
			books, err := bookshelf.List(store, test.query)
			testutils.AssertNoError(t, err)
			if !reflect.DeepEqual(books, test.want) {
				t.Errorf("got %v want %v", books, test.want)
			}
		})
	}
}

func TestUpdate(t *testing.T) {

	for _, test := range []struct {
		name     string
		fields   map[string]interface{}
		book     *bookshelf.Book
		wantBook *bookshelf.Book
	}{
		{
			"Update success: title only",
			map[string]interface{}{"title": "Romeo And Juliet"},
			testBooks[1],
			&bookshelf.Book{testBooks[1].ID, "Romeo And Juliet", testBooks[1].Author},
		},
		{
			"Update success: author only",
			map[string]interface{}{"author": "William Shakespeare"},
			testBooks[1],
			&bookshelf.Book{testBooks[1].ID, testBooks[1].Title, "William Shakespeare"},
		},
		{
			"Update success: title+author",
			map[string]interface{}{"title": "Romeo And Juliet", "author": "William Shakespeare"},
			testBooks[1],
			&bookshelf.Book{testBooks[1].ID, "Romeo And Juliet", "William Shakespeare"},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			store := testutils.NewSpyStore(testBooks)
			updated, err := bookshelf.Update(store, test.book.ID, test.fields)
			testutils.AssertNoError(t, err)
			testutils.AssertBooksEqual(t, updated, test.wantBook)
		})
	}

	for _, test := range []struct {
		name   string
		id     int64
		fields map[string]interface{}
		want   error
	}{
		{"Update failure: empty fields", 10, nil, bookshelf.ErrEmptyFields},
		{"Update failure: invalid fields", 10, map[string]interface{}{"isbn": 10}, bookshelf.ErrInvalidFields},
		{"Update failure: zero value ID", 0, map[string]interface{}{"author": "doesn't matter"}, bookshelf.ErrZeroValueID},
	} {
		t.Run(test.name, func(t *testing.T) {
			store := testutils.NewSpyStore(testBooks)
			_, err := bookshelf.Update(store, test.id, test.fields)
			testutils.AssertError(t, err, test.want)
		})
	}
}

func TestDelete(t *testing.T) {

	for _, test := range []struct {
		name string
		id   int64
		want *bookshelf.Book
	}{
		{"Delete success", 10, testBooks[0]},
	} {
		t.Run(test.name, func(t *testing.T) {
			store := testutils.NewSpyStore(testBooks)
			_, err := bookshelf.Delete(store, test.id)
			testutils.AssertNoError(t, err)
			var dummy bookshelf.Book
			err = store.ByID(&dummy, test.id)
			testutils.AssertError(t, err, bookshelf.ErrBookDoesNotExist)
		})
	}

	for _, test := range []struct {
		name string
		id   int64
		want error
	}{
		{"Delete failure: does not exist", 42, bookshelf.ErrBookDoesNotExist},
		{"Delete failure: zero value ID", 0, bookshelf.ErrZeroValueID},
	} {
		t.Run(test.name, func(t *testing.T) {
			store := testutils.NewSpyStore(testBooks)
			_, err := bookshelf.Delete(store, test.id)
			testutils.AssertError(t, err, test.want)
		})
	}
}

var (
	dummyBooks = make([]*bookshelf.Book, 0)
	testBooks  = []*bookshelf.Book{
		&bookshelf.Book{ID: 10, Author: "W. Shakespeare", Title: "The Tragedie of Hamlet"},
		&bookshelf.Book{ID: 22, Author: "W. Shakespeare", Title: "Romeo & Juliet"},
		&bookshelf.Book{ID: 24, Author: "Ernest Hemingway", Title: "The Old Man and The Sea"},
	}
)
