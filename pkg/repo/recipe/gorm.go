package recipe

import (
	"context"
	"fmt"

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
			Order:    i + 1,
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

func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(&Recipe{}, &RecipeStep{})
}

func (r *RepoGorm) FindByID(ctx context.Context, id int64) (*entity.Recipe, error) {
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

func (r *RepoGorm) Edit(ctx context.Context, recipe *entity.Recipe) error {
	rp := &Recipe{}
	rp.FromEntity(recipe)

	// Get actual steps
	var steps []*RecipeStep
	err := r.db.WithContext(ctx).Where("recipe_id = ?", rp.ID).Find(&steps).Error
	if err != nil {
		return err
	}
	//If there are less steps than before, delete the last ones (bulk delete). Also drop the ones deleted from the slice
	if len(steps) > len(rp.Steps) {
		err := r.db.WithContext(ctx).Delete(&RecipeStep{}, "recipe_id = ? AND `order` > ?", rp.ID, len(rp.Steps)).Error
		if err != nil {
			return err
		}
		steps = steps[:len(rp.Steps)]
	}
	// Asign id to steps that are not new
	for i, step := range steps {
		if i < len(steps) {
			step.Content = rp.Steps[i].Content
		}
	}
	rp.Steps = steps
	fmt.Println("recipe")
	fmt.Println(rp)
	fmt.Println("recipe")
	fmt.Println("steps")
	for _, step := range rp.Steps {
		fmt.Println(step)
	}
	fmt.Println("steps")
	// Save Recipe
	err = r.db.WithContext(ctx).Save(rp).Error
	if err != nil {
		return err
	}
	recipe.FromEntity(rp.ToEntity())
	return nil
}

func (r *RepoGorm) Delete(ctx context.Context, recipe *entity.Recipe) error {
	rp := &Recipe{}
	rp.FromEntity(recipe)
	// Delete steps
	err := r.db.WithContext(ctx).Delete(&RecipeStep{}, "recipe_id = ?", rp.ID).Error
	if err != nil {
		return err
	}
	// Delete recipe
	return r.db.WithContext(ctx).Delete(rp).Error
}
