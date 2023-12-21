package ingredient

import (
	"context"

	"gorm.io/gorm"

	"github.com/TomeuUris/recipes-catalog/pkg/entity"
)

type Ingredient struct {
	gorm.Model
	Name string
	Type string
}

func (i *Ingredient) ToEntity() *entity.Ingredient {
	return &entity.Ingredient{
		ID:   int64(i.ID),
		Name: i.Name,
		Type: i.Type,
	}
}

func (i *Ingredient) FromEntity(ingredient *entity.Ingredient) {
	i.ID = uint(ingredient.ID)
	i.Name = ingredient.Name
	i.Type = ingredient.Type
}

type RepoGorm struct {
	db *gorm.DB
}

func NewGormRepo(db *gorm.DB) *RepoGorm {
	return &RepoGorm{
		db: db,
	}
}

func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(&Ingredient{})
}

func (r *RepoGorm) FindByID(ctx context.Context, id int) (*entity.Ingredient, error) {
	ingredient := &Ingredient{}
	if err := r.db.WithContext(ctx).First(ingredient, id).Error; err != nil {
		return nil, err
	}

	return ingredient.ToEntity(), nil
}

func (r *RepoGorm) FindByFilter(ctx context.Context, f *FindFilter) ([]*entity.Ingredient, error) {
	var ingredients []*Ingredient
	if err := r.db.WithContext(ctx).Where("type = ?", f.Type).Find(&ingredients).Error; err != nil {
		return nil, err
	}
	result := make([]*entity.Ingredient, len(ingredients))
	for i, ingredient := range ingredients {
		result[i] = ingredient.ToEntity()
	}

	return result, nil
}

func (r *RepoGorm) Add(ctx context.Context, ingredient *entity.Ingredient) error {
	i := &Ingredient{}
	i.FromEntity(ingredient)
	err := r.db.WithContext(ctx).Create(i).Error
	if err != nil {
		return err
	}
	ingredient.ID = int64(i.ID)
	return nil
}

func (r *RepoGorm) Edit(ctx context.Context, ingredient *entity.Ingredient) error {
	i := &Ingredient{}
	i.FromEntity(ingredient)
	return r.db.WithContext(ctx).Save(i).Error
}

func (r *RepoGorm) Delete(ctx context.Context, ingredient *entity.Ingredient) error {
	i := &Ingredient{}
	i.FromEntity(ingredient)
	return r.db.WithContext(ctx).Delete(i).Error
}
