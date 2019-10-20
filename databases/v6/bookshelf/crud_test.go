package bookshelf_test

import (
	"testing"

	"github.com/djangulo/learn-go-with-tests/databases/v6/bookshelf"
	"github.com/djangulo/learn-go-with-tests/databases/v6/bookshelf/testutils"
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
				&bookshelf.Book{testBooks[1].ID + 1, "Moby Dick", "Herman Melville"},
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

var (
	dummyBooks = make([]*bookshelf.Book, 0)
	testBooks  = []*bookshelf.Book{
		&bookshelf.Book{ID: 10, Author: "W. Shakespeare", Title: "The Tragedie of Hamlet"},
		&bookshelf.Book{ID: 22, Author: "W. Shakespeare", Title: "Romeo & Juliet"},
	}
)
