package ingredient_test

import (
	"context"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/TomeuUris/recipes-catalog/pkg/entity"
	"github.com/TomeuUris/recipes-catalog/pkg/repo/ingredient"
)

var db *gorm.DB
var err error

func TestMain(m *testing.M) {
	// setup
	db, err = gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// Migrate the database schema
	err = db.AutoMigrate(&ingredient.Ingredient{})
	if err != nil {
		panic("failed to migrate database schema")
	}
	// run tests
	m.Run()
	// teardown
	db.Migrator().DropTable(&ingredient.Ingredient{})
}

func getExampleIngredientEntity() *entity.Ingredient {
	return &entity.Ingredient{
		Name: "ingredient",
		Type: "type",
	}
}

func getExampleIngredientGorm() *ingredient.Ingredient {
	return &ingredient.Ingredient{
		Name: "Spaghetti",
		Type: "Pasta",
	}
}

func TestRepoGorm_FindByFilter(t *testing.T) {
	tx := db.Begin()

	// Create a sample ingredient
	ingredientExample := getExampleIngredientGorm()
	if err := tx.Create(ingredientExample).Error; err != nil {
		t.Fatalf("failed to create ingredient: %v", err)
	}

	// Create the repo
	repo := ingredient.NewGormRepo(tx)

	// Find the ingredient
	ingredientsFound, err := repo.FindByFilter(context.Background(), &ingredient.FindFilter{
		Type: "Pasta",
	})
	if err != nil {
		t.Fatalf("failed to find ingredient: %v", err)
	}

	// Check number of ingredients found
	if len(ingredientsFound) != 1 {
		t.Fatalf("expected 1 ingredient, got %d", len(ingredientsFound))
	}

	// Check if ingredients found match the filter
	for _, ingredientFound := range ingredientsFound {
		if ingredientFound.Type != "Pasta" {
			t.Fatalf("expected ingredient type to be Pasta, got %s", ingredientFound.Type)
		}
	}
	tx.Rollback()
}

func TestRepoGorm_FindByID(t *testing.T) {
	tx := db.Begin()

	// Create a sample ingredient
	ingredientExample := getExampleIngredientGorm()
	if err := tx.Create(ingredientExample).Error; err != nil {
		t.Fatalf("failed to create ingredient: %v", err)
	}

	// Create the repo
	repo := ingredient.NewGormRepo(tx)

	// Find the ingredient
	ingredientFound, err := repo.FindByID(context.Background(), int(ingredientExample.ID))
	if err != nil {
		t.Fatalf("failed to find ingredient: %v", err)
	}

	// Check if ingredient found matches the ingredient added
	if ingredientFound.Name != ingredientExample.Name {
		t.Fatalf("expected ingredient name to be %s, got %s", ingredientExample.Name, ingredientFound.Name)
	}
	tx.Rollback()
}

func TestRepoGorm_Add(t *testing.T) {
	tx := db.Begin()

	// Create a sample ingredient
	ingredientExample := getExampleIngredientEntity()

	// Create the repo
	repo := ingredient.NewGormRepo(tx)

	// Add the ingredient
	err := repo.Add(context.Background(), ingredientExample)
	if err != nil {
		t.Fatalf("failed to add ingredient: %v", err)
	}

	// Check if ingredient was added
	ingredientFound, err := repo.FindByID(context.Background(), int(ingredientExample.ID))
	if err != nil {
		t.Fatalf("failed to find ingredient: %v", err)
	}

	// Check if ingredient found matches the ingredient added
	if ingredientFound.Name != ingredientExample.Name {
		t.Fatalf("expected ingredient name to be %s, got %s", ingredientExample.Name, ingredientFound.Name)
	}
	tx.Rollback()
}

func TestRepoGorm_Edit(t *testing.T) {
	tx := db.Begin()

	// Create a sample ingredient
	ingredientExample := getExampleIngredientGorm()
	if err := tx.Create(ingredientExample).Error; err != nil {
		t.Fatalf("failed to create ingredient: %v", err)
	}

	// Create the repo
	repo := ingredient.NewGormRepo(tx)

	// Edit the ingredient
	ingredientExample.Name = "Spaghetti Bolognese"
	err := repo.Edit(context.Background(), ingredientExample.ToEntity())
	if err != nil {
		t.Fatalf("failed to edit ingredient: %v", err)
	}

	// Check if ingredient was edited
	ingredientFound, err := repo.FindByID(context.Background(), int(ingredientExample.ID))
	if err != nil {
		t.Fatalf("failed to find ingredient: %v", err)
	}

	// Check if ingredient found matches the ingredient edited
	if ingredientFound.Name != ingredientExample.Name {
		t.Fatalf("expected ingredient name to be %s, got %s", ingredientExample.Name, ingredientFound.Name)
	}
	tx.Rollback()
}

func TestRepoGorm_Delete(t *testing.T) {
	tx := db.Begin()

	// Create a sample ingredient
	ingredientExample := getExampleIngredientGorm()
	if err := tx.Create(ingredientExample).Error; err != nil {
		t.Fatalf("failed to create ingredient: %v", err)
	}

	// Create the repo
	repo := ingredient.NewGormRepo(tx)

	// Delete the ingredient
	err := repo.Delete(context.Background(), ingredientExample.ToEntity())
	if err != nil {
		t.Fatalf("failed to delete ingredient: %v", err)
	}

	// Check if ingredient was deleted
	ingredientFound, err := repo.FindByID(context.Background(), int(ingredientExample.ID))
	if err == nil {
		t.Fatalf("expected to not find ingredient, got %v", ingredientFound)
	}
	tx.Rollback()
}


func TestRepoGorm_CountByFilter(t *testing.T) {
	tx := db.Begin()

	// Create a sample ingredient
	ingredientExample := getExampleIngredientGorm()
	if err := tx.Create(ingredientExample).Error; err != nil {
		t.Fatalf("failed to create ingredient: %v", err)
	}

	// Create the repo
	repo := ingredient.NewGormRepo(tx)

	// Find the ingredient
	ingredientsFound, err := repo.CountByFilter(&ingredient.FindFilter{
		Type: "Pasta",
	})
	if err != nil {
		t.Fatalf("failed to find ingredient: %v", err)
	}

	// Check number of ingredients found
	if ingredientsFound != 1 {
		t.Fatalf("expected 1 ingredient, got %d", ingredientsFound)
	}

	// Check with a different filter
	ingredientsFound, err = repo.CountByFilter(&ingredient.FindFilter{
		Type: "Pizza",
	})
	if err != nil {
		t.Fatalf("failed to find ingredient: %v", err)
	}

	// Check number of ingredients found
	if ingredientsFound != 0 {
		t.Fatalf("expected 0 ingredient, got %d", ingredientsFound)
	}
	tx.Rollback()
}