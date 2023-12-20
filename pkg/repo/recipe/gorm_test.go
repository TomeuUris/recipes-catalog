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
		return
	}

	// Assert the expected number of recipes
	expectedCount := 1 // set the expected count here
	if len(recipes) != expectedCount {
		t.Errorf("unexpected number of recipes, got: %d, want: %d", len(recipes), expectedCount)
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
	}

	// Create a new RepoGorm instance
	repo := recipe.NewGormRepo(tx)

	recipe, err := repo.FindByID(context.Background(), 1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	// Assert the recipe id is the expected one
	var expectedId int64 = 1 // set the expected id here
	if recipe.ID != expectedId {
		t.Errorf("unexpected recipe id, got: %d, want: %d", recipe.ID, expectedId)
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
		return
	}

	// Assert the recipe id is not zero
	if recipe.ID == 0 {
		t.Errorf("unexpected recipe id, got: %d, want not zero", recipe.ID)
		return
	}
	tx.Rollback()
}

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
