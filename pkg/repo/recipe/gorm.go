package recipe

import (
	"context"

	"github.com/TomeuUris/recipes-catalog/pkg/entity"
	repo "github.com/TomeuUris/recipes-catalog/pkg/repo/ingredient"
	"gorm.io/gorm"
)

type RecipeStep struct {
	gorm.Model
	Content  string
	Order    int  `gorm:"uniqueIndex:idx_recipe_order"`
	RecipeID uint `gorm:"uniqueIndex:idx_recipe_order"`
	Recipe   Recipe
}

type Recipe struct {
	gorm.Model
	Name        string
	Description string
	Ingredients []*repo.Ingredient `gorm:"many2many:recipe_ingredients;"`
	Steps       []*RecipeStep      `gorm:"constraint:OnDelete:CASCADE"`
}

func (r *Recipe) ToEntity() *entity.Recipe {
	return &entity.Recipe{
		ID:          int64(r.ID),
		Name:        r.Name,
		Description: r.Description,
		Ingredients: r.IngredientsToEntity(),
		Steps:       r.StepsToEntity(),
	}
}

func (r *Recipe) FromEntity(recipe *entity.Recipe) {
	if recipe.ID != 0 {
		r.ID = uint(recipe.ID)
	}
	r.Name = recipe.Name
	r.Description = recipe.Description
	r.IngredientsFromEntity(recipe.Ingredients)
	r.StepsFromEntity(recipe.Steps)
}

func (r *Recipe) IngredientsToEntity() []*entity.Ingredient {
	result := make([]*entity.Ingredient, len(r.Ingredients))
	for i, ingredient := range r.Ingredients {
		result[i] = ingredient.ToEntity()
	}
	return result
}

func (r *Recipe) IngredientsFromEntity(ingredients []*entity.Ingredient) {
	result := make([]*repo.Ingredient, len(ingredients))
	for i, ingredient := range ingredients {
		result[i] = &repo.Ingredient{}
		result[i].FromEntity(ingredient)
	}
	r.Ingredients = result
}

func (r *Recipe) StepsToEntity() []string {
	result := make([]string, len(r.Steps))
	for i, step := range r.Steps {
		result[i] = step.Content
	}
	return result
}

func (r *Recipe) StepsFromEntity(steps []string) {
	result := make([]*RecipeStep, len(steps))
	for i, step := range steps {
		result[i] = &RecipeStep{
			Content:  step,
			Order:    i,
			RecipeID: r.ID,
			Recipe:   *r,
		}
	}
	r.Steps = result
}

type RepoGorm struct {
	db *gorm.DB
}

func NewGormRepo(db *gorm.DB) *RepoGorm {
	return &RepoGorm{
		db: db,
	}
}

func (r *RepoGorm) FindByID(ctx context.Context, id int) (*entity.Recipe, error) {
	recipe := &Recipe{}
	if err := r.db.WithContext(ctx).
		Preload("Steps", func(db *gorm.DB) *gorm.DB {
			return db.Order("`recipe_steps`.`order` ASC")
		}).
		Preload("Ingredients").
		First(recipe, id).Error; err != nil {
		return nil, err
	}

	return recipe.ToEntity(), nil
}

type FindFilter struct {
	Id int
}

func (r *RepoGorm) FindByFilter(ctx context.Context, f *FindFilter) ([]*entity.Recipe, error) {
	var recipes []*Recipe
	if err := r.db.WithContext(ctx).
		Preload("Steps", func(db *gorm.DB) *gorm.DB {
			return db.Order("`recipe_steps`.`order` ASC")
		}).
		Preload("Ingredients").
		Where(f).
		Find(&recipes).
		Error; err != nil {
		return nil, err
	}

	result := make([]*entity.Recipe, len(recipes))
	for i, recipe := range recipes {
		result[i] = recipe.ToEntity()
	}

	return result, nil
}

func (r *RepoGorm) Add(ctx context.Context, recipe *entity.Recipe) error {
	rp := &Recipe{}
	rp.FromEntity(recipe)
	err := r.db.WithContext(ctx).Create(rp).Error
	if err != nil {
		return err
	}
	recipe.FromEntity(rp.ToEntity())
	return nil
}
