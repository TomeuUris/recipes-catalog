package payload

import "github.com/TomeuUris/recipes-catalog/pkg/entity"

type CookingUnit struct {
	Name *string `json:"name"`
}

type CookingUnits []CookingUnit

func (i *CookingUnit) ApplyTo(ingredient *entity.CookingUnit) {
	if i.Name != nil {
		ingredient.Name = *i.Name
	}
}

func (i *CookingUnit) ToEntity() *entity.CookingUnit {
	return &entity.CookingUnit{
		Name: *i.Name,
	}
}
