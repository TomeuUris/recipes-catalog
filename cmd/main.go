package main

import (
	"os"

	"github.com/TomeuUris/recipes-catalog/api/v1/controller"
	_ "github.com/TomeuUris/recipes-catalog/docs"
	"github.com/TomeuUris/recipes-catalog/pkg/repo/ingredient"
	"github.com/TomeuUris/recipes-catalog/pkg/repo/recipe"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Repo struct {
	Ingredients ingredient.Repo
	Recipes     recipe.Repo
}

func main() {
	db, err := OpenDB()
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	if err := RunMigrations(db); err != nil {
		panic(err)
	}

	ingredientsController := controller.NewIngredientController(ingredient.NewGormRepo(db))
	recipesController := controller.NewRecipeController(recipe.NewGormRepo(db))

	r := gin.Default()
	r = controller.SetupIngredientsRouter(ingredientsController, r)
	r = controller.SetupRecipesRouter(recipesController, r)
	if os.Getenv("ENV") != "prod" {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	r.Run()
}

func OpenDB() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open("database.sqlite"), &gorm.Config{})
}

func RunMigrations(db *gorm.DB) error {
	if err := ingredient.RunMigrations(db); err != nil {
		return err
	}
	if err := recipe.RunMigrations(db); err != nil {
		return err
	}
	return nil
}
