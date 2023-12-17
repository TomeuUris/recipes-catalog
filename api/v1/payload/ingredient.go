package payload

import "github.com/TomeuUris/recipes-catalog/pkg/entity"

type Ingredient struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Ingredients []Ingredient

func (i *Ingredient) ApplyTo(ingredient *entity.Ingredient) {
	ingredient.Name = i.Name
	ingredient.Type = i.Type
}

func (i *Ingredient) ToEntity() *entity.Ingredient {
	return &entity.Ingredient{
		Name: i.Name,
		Type: i.Type,
	}
}
