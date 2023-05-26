package ingredients

import (
	"github.com/alecthomas/assert/v2"
	"testing"
)

type StoreContract struct {
	NewStore func() Store
}

func (s StoreContract) Test(t *testing.T) {
	t.Run("it returns what is put in", func(t *testing.T) {
		want := []Ingredient{
			{Name: "Bananas", Quantity: 2},
			{Name: "Flour", Quantity: 1},
			{Name: "Eggs", Quantity: 2},
		}
		store := s.NewStore()
		store.Store(want...)

		got := store.GetIngredients()
		assert.Equal(t, got, want)
	})
}
