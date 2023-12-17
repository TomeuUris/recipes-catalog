package view

import "github.com/TomeuUris/recipes-catalog/pkg/entity"

type Recipe struct {
	ID          int64        `json:"id"`
	Name        string       `json:"name"`
	Ingredients []Ingredient `json:"ingredients"`
	Steps       []string     `json:"steps"`
}

func (r *Recipe) FromEntity(recipe *entity.Recipe) {
	r.ID = recipe.ID
	r.Name = recipe.Name
	r.Steps = recipe.Steps
	r.Ingredients = make([]Ingredient, len(recipe.Ingredients))
	for i, ingredient := range recipe.Ingredients {
		r.Ingredients[i].FromEntity(ingredient)
	}
}
