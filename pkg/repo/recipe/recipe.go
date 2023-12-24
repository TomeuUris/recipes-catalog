package recipe

import (
	"context"

	"github.com/TomeuUris/recipes-catalog/pkg/entity"
)

type Repo interface {
	FindByID(ctx context.Context, id int64) (*entity.Recipe, error)
	FindByFilter(ctx context.Context, f *FindFilter) ([]*entity.Recipe, error)
	CountByFilter(ctx context.Context, f *FindFilter) (int, error)
	Add(ctx context.Context, recipe *entity.Recipe) error
	Edit(ctx context.Context, recipe *entity.Recipe) error
	Delete(ctx context.Context, recipe *entity.Recipe) error
}

type FindFilter struct {
	Id int `form:"id"`
}
