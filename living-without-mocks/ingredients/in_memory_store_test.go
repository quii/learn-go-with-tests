package ingredients_test

import (
	"github.com/quii/learn-go-with-tests/living-without-mocks/ingredients"
	"testing"
)

func TestInMemoryStore(t *testing.T) {
	ingredients.StoreContract{
		NewStore: func() ingredients.Store {
			return &ingredients.InMemoryStore{}
		},
	}.Test(t)
}
