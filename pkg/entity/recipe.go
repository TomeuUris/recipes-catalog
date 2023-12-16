package entity

type Recipe struct {
	ID          int64
	Name        string
	Ingredients []*Ingredient
	Steps       []string
}

func NewRecipe(id int64, name string, ingredients []*Ingredient, steps []string) *Recipe {
	return &Recipe{
		ID:          id,
		Name:        name,
		Ingredients: ingredients,
		Steps:       steps,
	}
}
