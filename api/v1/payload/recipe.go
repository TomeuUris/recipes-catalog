package payload

import (
	"github.com/TomeuUris/recipes-catalog/api/v1/view"
	"github.com/TomeuUris/recipes-catalog/pkg/entity"
)

// Recipe is the payload for the recipe entity
type Recipe struct {
	Name        string            `json:"name"`
	Ingredients []view.Ingredient `json:"ingredients"`
	Steps       []string          `json:"steps"`
}

// Convert the payload to the entity
func (p *Recipe) ToEntity() *entity.Recipe {
	ingredients := make([]*entity.Ingredient, len(p.Ingredients))
	for i, ingredient := range p.Ingredients {
		ingredients[i] = &entity.Ingredient{
			ID:   ingredient.ID,
			Name: ingredient.Name,
			Type: ingredient.Type,
		}
	}
	return &entity.Recipe{
		Name:        p.Name,
		Ingredients: ingredients,
		Steps:       p.Steps,
	}
}
