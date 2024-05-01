package view

import "github.com/TomeuUris/recipes-catalog/pkg/entity"

type CookingUnit struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (i *CookingUnit) FromEntity(ingredient *entity.CookingUnit) {
	i.ID = ingredient.ID
	i.Name = ingredient.Name
}

func (i *CookingUnit) ToEntity() *entity.CookingUnit {
	return &entity.CookingUnit{
		ID:   i.ID,
		Name: i.Name,
	}
}
