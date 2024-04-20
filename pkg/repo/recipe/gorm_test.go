package recipe_test

import (
	"context"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/TomeuUris/recipes-catalog/pkg/entity"
	"github.com/TomeuUris/recipes-catalog/pkg/repo/ingredient"
	"github.com/TomeuUris/recipes-catalog/pkg/repo/recipe"
)

var db *gorm.DB
var err error
var ctx context.Context = context.Background()

func TestMain(m *testing.M) {
	// setup
	db, err = gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// Migrate the database schema
	err = db.AutoMigrate(&recipe.Recipe{}, &ingredient.Ingredient{}, &recipe.RecipeStep{})
	if err != nil {
		panic("failed to migrate database schema")
	}
	// run tests
	m.Run()
	// teardown
	db.Migrator().DropTable(&recipe.Recipe{}, &ingredient.Ingredient{}, &recipe.RecipeStep{})

}

func getExampleRecipeEntity() *entity.Recipe {
	return &entity.Recipe{
		Name:        "recipe",
		Description: "description",
		Ingredients: []*entity.Ingredient{
			{
				Name: "ingredient",
				Type: "type",
			},
		},
		Steps: []string{
			"step1",
			"step2",
		},
	}
}

func getExampleRecipeGorm() *recipe.Recipe {
	return &recipe.Recipe{
		Name:        "recipe",
		Description: "description",
		Ingredients: []*ingredient.Ingredient{
			{
				Name: "ingredient",
				Type: "type",
			},
		},
		Steps: []*recipe.RecipeStep{
			{
				Content: "step1",
				Order:   1,
			},
			{
				Content: "step2",
				Order:   2,
			},
		},
	}
}

func TestRepoGorm_FindByFilter(t *testing.T) {
	tx := db.Begin()

	// Create a sample recipe
	recipeExample := getExampleRecipeGorm()

	err = tx.Create(recipeExample).Error
	if err != nil {
		t.Fatalf("failed to create recipe: %v", err)
		tx.Rollback()
		return
	}

	// Create a new RepoGorm instance
	repo := recipe.NewGormRepo(tx)

	// Create a sample FindFilter
	filter := &recipe.FindFilter{
		Id: 1,
	}

	// Call the FindByFilter method
	recipes, err := repo.FindByFilter(context.Background(), filter)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		tx.Rollback()
		return
	}

	// Assert the expected number of recipes
	expectedCount := 1 // set the expected count here
	if len(recipes) != expectedCount {
		t.Errorf("unexpected number of recipes, got: %d, want: %d", len(recipes), expectedCount)
		tx.Rollback()
		return
	}

	tx.Rollback()
}

func TestRepoGorm_FindByID(t *testing.T) {
	tx := db.Begin()

	// Create a sample recipe
	expectedRecipe := getExampleRecipeGorm()

	err = tx.Create(expectedRecipe).Error
	if err != nil {
		t.Fatalf("failed to create recipe: %v", err)
		tx.Rollback()
	}

	// Create a new RepoGorm instance
	repo := recipe.NewGormRepo(tx)

	recipe, err := repo.FindByID(context.Background(), 1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		tx.Rollback()
		return
	}

	// Assert the recipe id is the expected one
	var expectedId int64 = 1 // set the expected id here
	if recipe.ID != expectedId {
		t.Errorf("unexpected recipe id, got: %d, want: %d", recipe.ID, expectedId)
		tx.Rollback()
		return
	}

	tx.Rollback()
}

func TestRepoGorm_Add(t *testing.T) {
	tx := db.Begin()

	// Create a new RepoGorm instance
	repo := recipe.NewGormRepo(tx)

	// Create a sample recipe
	recipe := getExampleRecipeEntity()

	// Call the Add method
	err = repo.Add(ctx, recipe)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		tx.Rollback()
		return
	}

	// Assert the recipe id is not zero
	if recipe.ID == 0 {
		t.Errorf("unexpected recipe id, got: %d, want not zero", recipe.ID)
		tx.Rollback()
		return
	}

	tx.Rollback()
}

func TestRepoGorm_AddDoesNotCreateMoreThanOneRecipe(t *testing.T) {
	tx := db.Begin()

	// Create a new RepoGorm instance
	repo := recipe.NewGormRepo(tx)

	recipes, err := repo.FindByFilter(ctx, &recipe.FindFilter{})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		tx.Rollback()
		return
	}

	// Assert the number of recipes is the expected one
	if len(recipes) != 0 {
		t.Errorf("unexpected number of recipes, got: %d, want: 0", len(recipes))
		tx.Rollback()
		return
	}

	// Create a sample recipe
	r := getExampleRecipeEntity()

	// Call the Add method
	err = repo.Add(ctx, r)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		tx.Rollback()
		return
	}

	newRecipes, err := repo.FindByFilter(ctx, &recipe.FindFilter{})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		tx.Rollback()
		return
	}

	// Assert the number of recipes is the expected one
	if len(newRecipes) != 1 {
		t.Errorf("unexpected number of recipes, got: %d, want: 1", len(newRecipes))
		tx.Rollback()
		return
	}

	tx.Rollback()
}

func TestRepoGorm_Edit_RemoveSteps(t *testing.T) {
	tx := db.Begin()

	// Create a sample recipe
	recipeExample := getExampleRecipeGorm()

	err = tx.Create(recipeExample).Error
	if err != nil {
		t.Fatalf("failed to create recipe: %v", err)
		tx.Rollback()
		return
	}

	// Create a new RepoGorm instance
	repo := recipe.NewGormRepo(tx)

	// Remove one step
	modRecipeExample := recipeExample.ToEntity()
	modRecipeExample.Steps = modRecipeExample.Steps[:len(modRecipeExample.Steps)-1]

	// Call the Edit method
	err = repo.Edit(ctx, modRecipeExample)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		tx.Rollback()
		return
	}

	recipeFound := &recipe.Recipe{}
	err = tx.Preload("Steps", func(db *gorm.DB) *gorm.DB {
		return db.Order("`recipe_steps`.`order` ASC")
	}).Preload("Ingredients").First(recipeFound, recipeExample.ID).Error
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		tx.Rollback()
		return
	}

	// Assert the number of steps is the expected one
	if len(recipeFound.Steps) != len(modRecipeExample.Steps) {
		t.Errorf("unexpected number of steps, got: %d, want: %d", len(recipeFound.Steps), len(modRecipeExample.Steps))
		tx.Rollback()
		return
	}

	tx.Rollback()
}

func TestRepoGorm_Edit_AddSteps(t *testing.T) {
	tx := db.Begin()

	// Create a sample recipe
	recipeExample := getExampleRecipeGorm()

	err = tx.Create(recipeExample).Error
	if err != nil {
		t.Fatalf("failed to create recipe: %v", err)
		tx.Rollback()
		return
	}

	// Create a new RepoGorm instance
	repo := recipe.NewGormRepo(tx)

	// Add one step
	modRecipeExample := recipeExample.ToEntity()
	modRecipeExample.Steps = append(modRecipeExample.Steps, "step3")

	// Call the Edit method
	err = repo.Edit(ctx, modRecipeExample)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		tx.Rollback()
		return
	}

	recipeFound := &recipe.Recipe{}
	err = tx.Preload("Steps", func(db *gorm.DB) *gorm.DB {
		return db.Order("`recipe_steps`.`order` ASC")
	}).Preload("Ingredients").First(recipeFound, recipeExample.ID).Error
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		tx.Rollback()
		return
	}

	// Assert the number of steps is the expected one
	if len(recipeFound.Steps) != len(modRecipeExample.Steps) {
		t.Errorf("unexpected number of steps, got: %d, want: %d", len(recipeFound.Steps), len(recipeExample.Steps))
		tx.Rollback()
		return
	}

	tx.Rollback()
}

func TestRepoGorm_Edit_ChangeValue(t *testing.T) {
	tx := db.Begin()

	// Create a sample recipe
	recipeExample := getExampleRecipeGorm()

	err = tx.Create(recipeExample).Error
	if err != nil {
		t.Fatalf("failed to create recipe: %v", err)
		tx.Rollback()
		return
	}

	// Create a new RepoGorm instance
	repo := recipe.NewGormRepo(tx)

	modRecipe := recipeExample.ToEntity()
	modRecipe.Name = "Spaghetti Bolognese"

	// Call the Edit method
	err = repo.Edit(ctx, modRecipe)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		tx.Rollback()
		return
	}

	recipeFound := &recipe.Recipe{}
	err = tx.First(recipeFound, modRecipe.ID).Error
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		tx.Rollback()
		return
	}

	// Assert the recipe name is the expected one
	if recipeFound.Name != modRecipe.Name {
		t.Errorf("unexpected recipe name, got: %s, want: %s", recipeFound.Name, modRecipe.Name)
		tx.Rollback()
		return
	}

	tx.Rollback()
}

func TestRepoGorm_Delete(t *testing.T) {
	tx := db.Begin()

	// Create a sample recipe
	recipeExample := getExampleRecipeGorm()

	err = tx.Create(recipeExample).Error
	if err != nil {
		t.Fatalf("failed to create recipe: %v", err)
	}

	// Create a new RepoGorm instance
	repo := recipe.NewGormRepo(tx)

	// Call the Delete method
	err = repo.Delete(ctx, recipeExample.ToEntity())
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	recipeFound := &recipe.Recipe{}
	err = tx.First(recipeFound, recipeExample.ID).Error
	if err == nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	// Check nested steps are deleted
	stepFound := &recipe.RecipeStep{}
	err = tx.First(stepFound, recipeExample.Steps[0].ID).Error
	if err == nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	tx.Rollback()
}

func TestRepoGorm_CountByFilter(t *testing.T) {
	tx := db.Begin()

	// Create a sample recipe
	recipeExample := getExampleRecipeGorm()

	err = tx.Create(recipeExample).Error
	if err != nil {
		t.Fatalf("failed to create recipe: %v", err)
	}

	// Create a new RepoGorm instance
	repo := recipe.NewGormRepo(tx)

	// Create a sample FindFilter
	filter := &recipe.FindFilter{
		Id: 1,
	}

	// Call the CountByFilter method
	count, err := repo.CountByFilter(ctx, filter)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	// Assert the expected number of recipes
	expectedCount := 1 // set the expected count here
	if count != expectedCount {
		t.Errorf("unexpected number of recipes, got: %d, want: %d", count, expectedCount)
		return
	}

	tx.Rollback()
}
