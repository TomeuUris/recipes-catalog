package recipe

import (
	"context"

	"github.com/TomeuUris/recipes-catalog/pkg/entity"
)

type Repo interface {
	FindByID(ctx context.Context, id int64) (*entity.Recipe, error)
	// FindByFilter(f *FindFilter) ([]*entity.Recipe, error)
	// CountByFilter(f *FindFilter) (int, error)
	Add(ctx context.Context, recipe *entity.Recipe) error
	// Edit(recipe *entity.Recipe) error
	// Delete(recipe *entity.Recipe) error
}

// type FindFilter struct {
// 	RecipeIDs []int64
// 	Name      string
// }
