package living_without_mocks

import "github.com/quii/learn-go-with-tests/living-without-mocks/ingredients"

type Recipe struct {
	Name        string
	Ingredients []ingredients.Ingredient
}

type RecipeBook interface {
	GetRecipes() []Recipe
}

type RecipeMatcher struct {
	recipeBook      RecipeBook
	ingredientStore ingredients.Store
}

func NewRecipeMatcher(recipes RecipeBook, ingredientStore ingredients.Store) *RecipeMatcher {
	return &RecipeMatcher{recipeBook: recipes, ingredientStore: ingredientStore}
}

func (m RecipeMatcher) SuggestRecipes() []Recipe {
	var suggestions []Recipe
	for _, recipe := range m.recipeBook.GetRecipes() {
		if m.canMake(recipe) {
			suggestions = append(suggestions, recipe)
		}
	}
	return suggestions
}

func (m RecipeMatcher) canMake(recipe Recipe) bool {
	for _, ingredient := range recipe.Ingredients {
		if !m.hasIngredient(ingredient) {
			return false
		}
	}
	return true
}

func (m RecipeMatcher) hasIngredient(ingredient ingredients.Ingredient) bool {
	for _, pantryIngredient := range m.ingredientStore.GetIngredients() {
		if pantryIngredient.Name == ingredient.Name {
			return true
		}
	}
	return false
}
