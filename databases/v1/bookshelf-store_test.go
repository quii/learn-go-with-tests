package bookshelf_test

import (
	"github.com/google/go-cmp/cmp"
	bookshelf "github.com/quii/learn-go-with-tests/databases/v1"
	"testing"
)

func TestBookShelfStore(t *testing.T) {
	store := bookshelf.NewStore()
	t.Run("it stores a book", func(t *testing.T) {
		book := bookshelf.Book{
			Title:  "East of Eden",
			Author: "John Steinbeck",
		}
		store.StoreBook(book)
		fetchedBooks, err := store.GetBooks()

		if err != nil {
			t.Fatal(err)
		}

		if !cmp.Equal(fetchedBooks[0], book) {
			t.Errorf("got %+v wanted %+v", fetchedBooks[0], book)
		}
	})
}
