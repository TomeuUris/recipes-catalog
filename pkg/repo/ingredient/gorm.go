package ingredient

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/TomeuUris/recipes-catalog/pkg/entity"
)

// Database model
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

// Repository implementation
type RepoGorm struct {
	db *gorm.DB
}

// Utility functions
func NewGormRepo(db *gorm.DB) *RepoGorm {
	return &RepoGorm{
		db: db,
	}
}

func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(&Ingredient{})
}

// CRUD functions
func (r *RepoGorm) FindByID(ctx context.Context, id int) (*entity.Ingredient, error) {
	ingredient := &Ingredient{}
	if err := r.db.WithContext(ctx).First(ingredient, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entity.ErrNotFound
		}
		return nil, err
	}

	return ingredient.ToEntity(), nil
}

func (r *RepoGorm) FindByFilter(ctx context.Context, f *FindFilter) ([]*entity.Ingredient, error) {
	var ingredients []*Ingredient
	if err := r.db.WithContext(ctx).Where(f).Find(&ingredients).Error; err != nil {
		return nil, err
	}
	result := make([]*entity.Ingredient, len(ingredients))
	for i, ingredient := range ingredients {
		result[i] = ingredient.ToEntity()
	}

	return result, nil
}

func (r *RepoGorm) CountByFilter(f *FindFilter) (int, error) {
	var count int64
	if err := r.db.Model(&Ingredient{}).Where(f).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
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
	err := r.db.WithContext(ctx).Save(i).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.ErrNotFound
		}
		return err
	}
	return nil
}

func (r *RepoGorm) Delete(ctx context.Context, ingredient *entity.Ingredient) error {
	i := &Ingredient{}
	i.FromEntity(ingredient)
	err := r.db.WithContext(ctx).Delete(i).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.ErrNotFound
		}
		return err
	}
	return nil
}
