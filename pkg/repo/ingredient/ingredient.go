package ingredient

import (
	"context"

	"github.com/TomeuUris/recipes-catalog/pkg/entity"
)

type Repo interface {
	FindByID(ctx context.Context, id int) (*entity.Ingredient, error)
	FindByFilter(ctx context.Context, f *FindFilter) ([]*entity.Ingredient, error)
	// CountByFilter(f *FindFilter) (int, error)
	Add(ctx context.Context, ingredient *entity.Ingredient) error
	// Edit(recipe *entity.Recipe) error
	// Delete(recipe *entity.Recipe) error
}

type FindFilter struct {
	RecipeId int
}
