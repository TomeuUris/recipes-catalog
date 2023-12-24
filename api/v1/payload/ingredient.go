package payload

import "github.com/TomeuUris/recipes-catalog/pkg/entity"

type Ingredient struct {
	Name *string `json:"name"`
	Type *string `json:"type"`
}

type Ingredients []Ingredient

func (i *Ingredient) ApplyTo(ingredient *entity.Ingredient) {
	if i.Name != nil {
		ingredient.Name = *i.Name
	}
	if i.Type != nil {
		ingredient.Type = *i.Type
	}
}

func (i *Ingredient) ToEntity() *entity.Ingredient {
	return &entity.Ingredient{
		Name: *i.Name,
		Type: *i.Type,
	}
}
