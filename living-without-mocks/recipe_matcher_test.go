package living_without_mocks

import (
	"github.com/quii/learn-go-with-tests/living-without-mocks/ingredients"
	"testing"
)
import "github.com/alecthomas/assert/v2"

var (
	bananaBread = Recipe{
		Name: "Banana Bread",
		Ingredients: []ingredients.Ingredient{
			{Name: "Bananas", Quantity: 2},
			{Name: "Flour", Quantity: 1},
			{Name: "Eggs", Quantity: 2},
		},
	}
	bananaMilkshake = Recipe{
		Name: "Banana Milkshake",
		Ingredients: []ingredients.Ingredient{
			{Name: "Bananas", Quantity: 2},
			{Name: "Milk", Quantity: 1},
		},
	}
	recipeStore = InMemoryRecipeStore{Recipes: []Recipe{bananaBread, bananaMilkshake}}
)

func TestRecipeMatcher(t *testing.T) {

	t.Run("if we have no ingredients we can't make anything", func(t *testing.T) {
		assertAvailableRecipes(t, &ingredients.InMemoryStore{}, []Recipe{})
	})

	t.Run("if we have the ingredients for banana bread we can make it", func(t *testing.T) {
		store := &ingredients.InMemoryStore{}
		store.Store(
			ingredients.Ingredient{Name: "Bananas", Quantity: 2},
			ingredients.Ingredient{Name: "Flour", Quantity: 1},
			ingredients.Ingredient{Name: "Eggs", Quantity: 2},
		)
		assertAvailableRecipes(t, store, []Recipe{bananaBread})
	})

	t.Run("if we have bananas and milk, we can make banana milkshake", func(t *testing.T) {
		store := &ingredients.InMemoryStore{}
		store.Store(
			ingredients.Ingredient{Name: "Bananas", Quantity: 2},
			ingredients.Ingredient{Name: "Milk", Quantity: 1},
		)
		assertAvailableRecipes(t, store, []Recipe{bananaMilkshake})
	})

	t.Run("if we have ingredients for banana bread and milkshake, we can make both", func(t *testing.T) {
		store := &ingredients.InMemoryStore{}
		store.Store(
			ingredients.Ingredient{Name: "Bananas", Quantity: 2},
			ingredients.Ingredient{Name: "Flour", Quantity: 1},
			ingredients.Ingredient{Name: "Eggs", Quantity: 2},
			ingredients.Ingredient{Name: "Milk", Quantity: 1},
		)
		assertAvailableRecipes(t, store, []Recipe{bananaMilkshake, bananaBread})
	})

}

func assertAvailableRecipes(t *testing.T, ingredientStore ingredients.Store, expectedRecipes []Recipe) {
	suggestions := NewRecipeMatcher(recipeStore, ingredientStore).SuggestRecipes()

	// create a map to count occurrences of each recipe in the suggestions
	suggestionCounts := make(map[string]int)
	for _, suggestion := range suggestions {
		suggestionCounts[suggestion.Name]++
	}

	// check that the counts of the expected recipes match the actual counts in the suggestions
	for _, expectedRecipe := range expectedRecipes {
		actualCount, ok := suggestionCounts[expectedRecipe.Name]
		if !ok {
			t.Errorf("expected recipe %s not found in suggestions", expectedRecipe.Name)
			continue
		}
		if actualCount != 1 {
			t.Errorf("expected recipe %s to appear once in suggestions, but found %d occurrences", expectedRecipe.Name, actualCount)
		}
	}
	// check that the number of suggestions matches the expected number of recipes
	assert.Equal(t, len(suggestions), len(expectedRecipes))
}
