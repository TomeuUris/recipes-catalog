package cooking_unit_test

import (
	"context"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/TomeuUris/recipes-catalog/pkg/entity"
	"github.com/TomeuUris/recipes-catalog/pkg/repo/cooking_unit"
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
	err = db.AutoMigrate(&cooking_unit.CookingUnit{})
	if err != nil {
		panic("failed to migrate database schema")
	}
	// run tests
	m.Run()
	// teardown
	db.Migrator().DropTable(&cooking_unit.CookingUnit{})
}

func getExampleCookingUnitEntity() *entity.CookingUnit {
	return &entity.CookingUnit{
		Name: "gram",
	}
}

func getExampleCookingUnitGorm() *cooking_unit.CookingUnit {
	return &cooking_unit.CookingUnit{
		Name: "gram",
	}
}

func TestRepoGorm_FindByFilter(t *testing.T) {
	tx := db.Begin()

	// Create a sample ingredient
	cookingUnitExample := getExampleCookingUnitGorm()
	if err := tx.Create(cookingUnitExample).Error; err != nil {
		t.Fatalf("failed to create cooking unit: %v", err)
	}

	// Create the repo
	repo := cooking_unit.NewGormRepo(tx)

	// Find the ingredient
	cookingUnitsFound, err := repo.FindByFilter(context.Background(), &cooking_unit.FindFilter{
		Name: "gram",
	})
	if err != nil {
		t.Fatalf("failed to find cooking unit: %v", err)
	}

	// Check number of ingredients found
	if len(cookingUnitsFound) != 1 {
		t.Fatalf("expected 1 ingredient, got %d", len(cookingUnitsFound))
	}

	// Check if ingredients found match the filter
	for _, cookingUnitFound := range cookingUnitsFound {
		if cookingUnitFound.Name != "gram" {
			t.Fatalf("expected cooking unit name to be gram, got %s", cookingUnitFound.Name)
		}
	}
	tx.Rollback()
}

func TestRepoGorm_FindByID(t *testing.T) {
	tx := db.Begin()

	// Create a sample ingredient
	cookingUnitExample := getExampleCookingUnitGorm()
	if err := tx.Create(cookingUnitExample).Error; err != nil {
		t.Fatalf("failed to create ingredient: %v", err)
	}

	// Create the repo
	repo := cooking_unit.NewGormRepo(tx)

	// Find the ingredient
	cookingUnitFound, err := repo.FindByID(context.Background(), int(cookingUnitExample.ID))
	if err != nil {
		t.Fatalf("failed to find ingredient: %v", err)
	}

	// Check if ingredient found matches the ingredient added
	if cookingUnitFound.Name != cookingUnitExample.Name {
		t.Fatalf("expected ingredient name to be %s, got %s", cookingUnitExample.Name, cookingUnitFound.Name)
	}
	tx.Rollback()
}

func TestRepoGorm_Add(t *testing.T) {
	tx := db.Begin()

	// Create a sample ingredient
	cookingUnitExample := getExampleCookingUnitEntity()

	// Create the repo
	repo := cooking_unit.NewGormRepo(tx)

	// Add the ingredient
	err := repo.Add(context.Background(), cookingUnitExample)
	if err != nil {
		t.Fatalf("failed to add ingredient: %v", err)
	}

	// Check if ingredient was added
	cookingUnitFound, err := repo.FindByID(context.Background(), int(cookingUnitExample.ID))
	if err != nil {
		t.Fatalf("failed to find ingredient: %v", err)
	}

	// Check if ingredient found matches the ingredient added
	if cookingUnitFound.Name != cookingUnitExample.Name {
		t.Fatalf("expected ingredient name to be %s, got %s", cookingUnitExample.Name, cookingUnitFound.Name)
	}
	tx.Rollback()
}

func TestRepoGorm_Edit(t *testing.T) {
	tx := db.Begin()

	// Create a sample ingredient
	cookingUnitExample := getExampleCookingUnitGorm()
	if err := tx.Create(cookingUnitExample).Error; err != nil {
		t.Fatalf("failed to create ingredient: %v", err)
	}

	// Create the repo
	repo := cooking_unit.NewGormRepo(tx)

	// Edit the ingredient
	cookingUnitExample.Name = "Spaghetti Bolognese"
	err := repo.Edit(context.Background(), cookingUnitExample.ToEntity())
	if err != nil {
		t.Fatalf("failed to edit ingredient: %v", err)
	}

	// Check if ingredient was edited
	cookingUnitFound, err := repo.FindByID(context.Background(), int(cookingUnitExample.ID))
	if err != nil {
		t.Fatalf("failed to find ingredient: %v", err)
	}

	// Check if ingredient found matches the ingredient edited
	if cookingUnitFound.Name != cookingUnitExample.Name {
		t.Fatalf("expected ingredient name to be %s, got %s", cookingUnitExample.Name, cookingUnitFound.Name)
	}
	tx.Rollback()
}

func TestRepoGorm_Delete(t *testing.T) {
	tx := db.Begin()

	// Create a sample ingredient
	cookingUnitExample := getExampleCookingUnitGorm()
	if err := tx.Create(cookingUnitExample).Error; err != nil {
		t.Fatalf("failed to create ingredient: %v", err)
	}

	// Create the repo
	repo := cooking_unit.NewGormRepo(tx)

	// Delete the ingredient
	err := repo.Delete(context.Background(), cookingUnitExample.ToEntity())
	if err != nil {
		t.Fatalf("failed to delete ingredient: %v", err)
	}

	// Check if ingredient was deleted
	cookingUnitFound, err := repo.FindByID(context.Background(), int(cookingUnitExample.ID))
	if err == nil {
		t.Fatalf("expected to not find ingredient, got %v", cookingUnitFound)
	}
	tx.Rollback()
}

func TestRepoGorm_CountByFilter(t *testing.T) {
	tx := db.Begin()

	// Create a sample ingredient
	cookingUnitExample := getExampleCookingUnitGorm()
	if err := tx.Create(cookingUnitExample).Error; err != nil {
		t.Fatalf("failed to create ingredient: %v", err)
	}

	// Create the repo
	repo := cooking_unit.NewGormRepo(tx)

	// Find the ingredient
	cookingUnitFound, err := repo.CountByFilter(&cooking_unit.FindFilter{
		Name: "gram",
	})
	if err != nil {
		t.Fatalf("failed to find ingredient: %v", err)
	}

	// Check number of ingredients found
	if cookingUnitFound != 1 {
		t.Fatalf("expected 1 ingredient, got %d", cookingUnitFound)
	}

	// Check with a different filter
	cookingUnitFound, err = repo.CountByFilter(&cooking_unit.FindFilter{
		Name: "kilogram",
	})
	if err != nil {
		t.Fatalf("failed to find ingredient: %v", err)
	}

	// Check number of ingredients found
	if cookingUnitFound != 0 {
		t.Fatalf("expected 0 ingredient, got %d", cookingUnitFound)
	}
	tx.Rollback()
}
