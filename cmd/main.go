package main

import (
	"os"

	"github.com/TomeuUris/recipes-catalog/api/v1/controller"
	_ "github.com/TomeuUris/recipes-catalog/docs"
	"github.com/TomeuUris/recipes-catalog/pkg/repo/ingredient"
	"github.com/TomeuUris/recipes-catalog/pkg/repo/recipe"
	"github.com/TomeuUris/recipes-catalog/pkg/sqlite"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Repo struct {
	Ingredients ingredient.Repo
	Recipes     recipe.Repo
}

func main() {
	db := sqlite.MustOpen("database.sqlite")
	ingredientsController := controller.NewIngredientController(ingredient.NewRepo(db))
	recipesController := controller.NewRecipeController(recipe.NewRepo(db))

	r := gin.Default()
	r = controller.SetupIngredientsRouter(ingredientsController, r)
	r = controller.SetupRecipesRouter(recipesController, r)
	if os.Getenv("ENV") != "prod" {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	r.Run()
}

// func test() {
// 	db := sqlite.MustOpen("database.sqlite")
// 	ctx := context.Background()

// 	repos := Repo{
// 		Ingredients: ingredient.NewRepo(db),
// 		Recipes:     recipe.NewRepo(db),
// 	}

// 	espaguetis := &entity.Ingredient{
// 		Name: "spaghetti",
// 		Type: "pasta",
// 	}

// 	tomate := &entity.Ingredient{
// 		Name: "tomate",
// 		Type: "salsa",
// 	}

// 	repos.Ingredients.Add(ctx, espaguetis)
// 	fmt.Println(espaguetis.ID)
// 	repos.Ingredients.Add(ctx, tomate)
// 	fmt.Println(tomate.ID)

// 	espaguetisConTomate := &entity.Recipe{
// 		Name:        "spaghetti con tomate",
// 		Ingredients: []*entity.Ingredient{espaguetis, tomate},
// 		Steps: []string{
// 			"Cocer la pasta",
// 			"Preparar la salsa",
// 			"Mezclar la pasta con la salsa",
// 		},
// 	}

// 	repos.Recipes.Add(ctx, espaguetisConTomate)
// 	fmt.Println("ID: ", espaguetisConTomate.ID)

// 	recipe, err := repos.Recipes.FindByID(ctx, espaguetisConTomate.ID)
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println("ID: ", recipe.ID)
// 	fmt.Println("Name: ", recipe.Name)
// 	fmt.Println("Ingredients: ")
// 	for _, ingredient := range recipe.Ingredients {
// 		fmt.Println(" - ", ingredient.Name, ingredient.Type)
// 	}
// 	fmt.Println("Steps: ")
// 	for i, step := range recipe.Steps {
// 		fmt.Println(" ", i+1, " ", step)
// 	}
// }
