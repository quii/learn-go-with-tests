package living_without_mocks

import "testing"
import "github.com/alecthomas/assert/v2"

var (
	bananaBread = Recipe{
		Name: "Banana Bread",
		Ingredients: []Ingredient{
			{Name: "Bananas", Quantity: 2},
			{Name: "Flour", Quantity: 1},
			{Name: "Eggs", Quantity: 2},
		},
	}
	bananaMilkshake = Recipe{
		Name: "Banana Milkshake",
		Ingredients: []Ingredient{
			{Name: "Bananas", Quantity: 2},
			{Name: "Milk", Quantity: 1},
		},
	}
	recipeStore = InMemoryRecipeStore{Recipes: []Recipe{bananaBread, bananaMilkshake}}
)

func TestRecipeMatcher(t *testing.T) {
	assertAvailableRecipes := func(t *testing.T, ingredients []Ingredient, expectedRecipes []Recipe) {
		ingredientsStore := InMemoryIngredientStore{Ingredients: ingredients}
		recipeMatcher := NewRecipeMatcher(recipeStore, ingredientsStore)
		suggestions := recipeMatcher.SuggestRecipes()

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

	t.Run("if we have no ingredients we can't make anything", func(t *testing.T) {
		assertAvailableRecipes(t, []Ingredient{}, []Recipe{})
	})

	t.Run("if we have the ingredients for banana bread we can make it", func(t *testing.T) {
		assertAvailableRecipes(t, []Ingredient{
			{Name: "Bananas", Quantity: 2},
			{Name: "Flour", Quantity: 1},
			{Name: "Eggs", Quantity: 2},
		}, []Recipe{bananaBread})
	})

	t.Run("if we have bananas and milk, we can make banana milkshake", func(t *testing.T) {
		assertAvailableRecipes(t, []Ingredient{
			{Name: "Bananas", Quantity: 2},
			{Name: "Milk", Quantity: 1},
		}, []Recipe{bananaMilkshake})
	})

	t.Run("if we have ingredients for banana bread and milkshake, we can make both", func(t *testing.T) {
		assertAvailableRecipes(t, []Ingredient{
			{Name: "Bananas", Quantity: 2},
			{Name: "Flour", Quantity: 1},
			{Name: "Eggs", Quantity: 2},
			{Name: "Milk", Quantity: 1},
		}, []Recipe{bananaMilkshake, bananaBread})
	})

}
