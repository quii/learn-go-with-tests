package living_without_mocks

type InMemoryRecipeStore struct {
	Recipes []Recipe
}

func (s InMemoryRecipeStore) GetRecipes() []Recipe {
	return s.Recipes
}

type InMemoryIngredientStore struct {
	Ingredients []Ingredient
}

func (s InMemoryIngredientStore) GetIngredients() []Ingredient {
	return s.Ingredients
}
