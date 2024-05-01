package cooking_unit

import (
	"context"

	"github.com/TomeuUris/recipes-catalog/pkg/entity"
)

type Repo interface {
	FindByID(ctx context.Context, id int) (*entity.CookingUnit, error)
	FindByFilter(ctx context.Context, f *FindFilter) ([]*entity.CookingUnit, error)
	CountByFilter(f *FindFilter) (int, error)
	Add(ctx context.Context, unit *entity.CookingUnit) error
	Edit(ctx context.Context, unit *entity.CookingUnit) error
	Delete(ctx context.Context, unit *entity.CookingUnit) error
}

type FindFilter struct {
	Name string `form:"name"`
}
