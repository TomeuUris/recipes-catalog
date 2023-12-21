package view

import "github.com/TomeuUris/recipes-catalog/pkg/entity"

type Ingredient struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

func (i *Ingredient) FromEntity(ingredient *entity.Ingredient) {
	i.ID = ingredient.ID
	i.Name = ingredient.Name
	i.Type = ingredient.Type
}

func (i *Ingredient) ToEntity() *entity.Ingredient {
	return &entity.Ingredient{
		ID:   i.ID,
		Name: i.Name,
		Type: i.Type,
	}
}
