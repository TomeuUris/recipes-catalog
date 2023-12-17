package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/TomeuUris/recipes-catalog/api/v1/payload"
	"github.com/TomeuUris/recipes-catalog/api/v1/view"
	"github.com/TomeuUris/recipes-catalog/pkg/entity"
	"github.com/TomeuUris/recipes-catalog/pkg/repo/recipe"
	"github.com/gin-gonic/gin"
)

type RecipeController struct {
	repo recipe.Repo
}

func NewRecipeController(repo recipe.Repo) *RecipeController {
	return &RecipeController{repo: repo}
}

// @Summary Get recipe by ID
// @Description get User by its ID
// @Tags recipes
// @Produce  json
// @Param   id     path    int     true        "recipe ID"
// @Success 200 {object} view.Recipe
// @Router /recipes/{id} [get]
func (c *RecipeController) GetRecipeByIdHandler(ctx *gin.Context) {
	// Get the recipe ID from the URL parameter
	recipeIDStr := ctx.Params.ByName("id")
	recipeID, err := strconv.ParseInt(recipeIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recipe ID"})
		return
	}

	// Get the recipe from the database
	// dbctx := ctx.Request.Context()
	recipe, err := c.repo.FindByID(ctx, recipeID)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	recipeView := &view.Recipe{}
	recipeView.FromEntity(recipe)

	// Return the recipes as a response
	ctx.JSON(http.StatusOK, recipeView)
}

// @Summary Create recipe
// @Description Create a new recipe
// @Tags recipes
// @Accept  json
// @Produce  json
// @Param   recipe     body    payload.Recipe     true        "Recipe info"
// @Success 200 {object} view.Recipe
// @Router /recipes [post]
func (c *RecipeController) CreateRecipeHandler(ctx *gin.Context) {
	// Parse the request payload
	var recipePayload payload.Recipe
	if err := ctx.ShouldBindJSON(&recipePayload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	recipe := recipePayload.ToEntity()
	// Logic to create the recipe in the database
	err := c.repo.Add(ctx, recipe)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	recipeView := &view.Recipe{}
	recipeView.FromEntity(recipe)

	// Return the created recipe as a response
	ctx.JSON(http.StatusCreated, recipeView)
}

func SetupRecipesRouter(controller *RecipeController, router *gin.Engine) *gin.Engine {
	router.GET("/recipes/:id", controller.GetRecipeByIdHandler)
	router.POST("/recipes", controller.CreateRecipeHandler)
	return router
}
