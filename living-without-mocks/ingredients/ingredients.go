package ingredients

type Ingredient struct {
	Name     string
	Quantity int
}

type Store interface {
	GetIngredients() []Ingredient
	Store(...Ingredient)
}
