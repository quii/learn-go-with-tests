package living_without_mocks

type Recipe struct {
	Name        string
	Ingredients []Ingredient
}

type Ingredient struct {
	Name     string
	Quantity int
}

type RecipeBook interface {
	GetRecipes() []Recipe
}

type IngredientStore interface {
	GetIngredients() []Ingredient
}
type RecipeMatcher struct {
	recipeBook      RecipeBook
	ingredientStore IngredientStore
}

func NewRecipeMatcher(recipes RecipeBook, ingredientStore IngredientStore) *RecipeMatcher {
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

func (m RecipeMatcher) hasIngredient(ingredient Ingredient) bool {
	for _, pantryIngredient := range m.ingredientStore.GetIngredients() {
		if pantryIngredient.Name == ingredient.Name {
			return true
		}
	}
	return false
}
