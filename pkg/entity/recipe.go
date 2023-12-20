package entity

type Recipe struct {
	ID          int64
	Name        string
	Description string
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

func (r *Recipe) FromEntity(recipe *Recipe) error {
	r.ID = recipe.ID
	r.Name = recipe.Name
	r.Description = recipe.Description
	r.Ingredients = recipe.Ingredients
	r.Steps = recipe.Steps
	return nil
}
