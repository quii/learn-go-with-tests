package living_without_mocks

type InMemoryRecipeStore struct {
	Recipes []Recipe
}

func (s InMemoryRecipeStore) GetRecipes() []Recipe {
	return s.Recipes
}
