package main

import (
	_ "github.com/TomeuUris/recipes-catalog/docs"
)

func mainSqlite() {
	// db := sqlite.MustOpen("database.sqlite")
	// ingredientsController := controller.NewIngredientController(ingredient.NewGoRepo(db))
	// recipesController := controller.NewRecipeController(recipe.NewRepo(db))

	// r := gin.Default()
	// r = controller.SetupIngredientsRouter(ingredientsController, r)
	// r = controller.SetupRecipesRouter(recipesController, r)
	// if os.Getenv("ENV") != "prod" {
	// 	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// }

	// r.Run()
}
