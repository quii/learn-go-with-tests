package bookshelf_test

import (
	"testing"

	"github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf"
	"github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf/testutils"
)

func TestCreate(t *testing.T) {

	t.Run("can create a book", func(t *testing.T) {
		store := testutils.NewSpyStore()

		var book bookshelf.Book
		err := store.CreateBook(&book, "Moby Dick", "Herman Melville")
		testutils.AssertNoError(t, err)

		if book.ID == 0 {
			t.Error("book returned without an ID")
		}
	})

}
