package ingredients

type InMemoryStore struct {
	ingredients []Ingredient
}

func (s *InMemoryStore) Store(i ...Ingredient) {
	s.ingredients = append(s.ingredients, i...)
}

func (s *InMemoryStore) GetIngredients() []Ingredient {
	return s.ingredients
}
