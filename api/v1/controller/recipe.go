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

// @Summary Get recipes by filter
// @Description Retrieve recipes by filter
// @Tags recipes
// @Produce  json
// @Param   filter     query    recipe.FindFilter     true        "Filter parameters"
// @Success 200 {array} view.Recipe
// @Router /recipes [get]
func (c *RecipeController) GetRecipesByFilterHandler(ctx *gin.Context) {
	// Parse the request query
	var filter recipe.FindFilter
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the recipes from the database
	recipes, err := c.repo.FindByFilter(ctx, &filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert the recipes to the view model
	recipesView := make([]*view.Recipe, len(recipes))
	for i, recipe := range recipes {
		recipeView := &view.Recipe{}
		recipeView.FromEntity(recipe)
		recipesView[i] = recipeView
	}

	// Return the recipes as a response
	ctx.JSON(http.StatusOK, recipesView)
}

// @Summary Get recipe by ID
// @Description Retrieves a recipe by ID
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

// @Summary Count recipes by filter
// @Description Retrieves the number of recipes filtered by the given parameters
// @Tags recipes
// @Produce  json
// @Param   filter     query    recipe.FindFilter     true        "Filter parameters"
// @Success 200 {object} int
// @Router /recipes/count [get]
func (c *RecipeController) CountRecipeByFilterHandler(ctx *gin.Context) {
	// Parse the filter from the query parameters
	var filter recipe.FindFilter
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Count the recipes in the database
	count, err := c.repo.CountByFilter(ctx, &filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the recipes as a response
	ctx.JSON(http.StatusOK, count)
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

// @Summary Edit recipe
// @Description Edit an existing recipe
// @Tags recipes
// @Accept  json
// @Produce  json
// @Param   id     path    int     true        "recipe ID"
// @Param   recipe     body    payload.Recipe     true        "Recipe info"
// @Success 200 {object} view.Recipe
// @Router /recipes/{id} [patch]
func (c *RecipeController) EditRecipeHandler(ctx *gin.Context) {
	// Get the recipe ID from the URL parameter
	recipeIDStr := ctx.Params.ByName("id")
	recipeID, err := strconv.ParseInt(recipeIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recipe ID"})
		return
	}

	// Parse the request payload
	var recipePayload payload.Recipe
	if err := ctx.ShouldBindJSON(&recipePayload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the recipe from the database
	targetRecipe, err := c.repo.FindByID(ctx, recipeID)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update the recipe entity
	recipePayload.ApplyTo(targetRecipe)

	// Update the recipe in the database
	err = c.repo.Edit(ctx, targetRecipe)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert the recipe to a view
	recipeView := &view.Recipe{}
	recipeView.FromEntity(targetRecipe)

	// Return the updated recipe as a response
	ctx.JSON(http.StatusOK, recipeView)
}

// @Summary Delete recipe
// @Description Delete an existing recipe
// @Tags recipes
// @Accept  json
// @Produce  json
// @Param   id     path    int     true        "recipe ID"
// @Success 200
// @Router /recipes/{id} [delete]
func (c *RecipeController) DeleteRecipeHandler(ctx *gin.Context) {
	// Get the recipe ID from the URL parameter
	recipeIDStr := ctx.Params.ByName("id")
	recipeID, err := strconv.ParseInt(recipeIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recipe ID"})
		return
	}

	// Get the recipe from the database
	targetRecipe, err := c.repo.FindByID(ctx, recipeID)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Delete the recipe from the database
	err = c.repo.Delete(ctx, targetRecipe)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return a response with status 200 OK
	ctx.JSON(http.StatusOK, nil)
}

// SetupRecipesRouter sets up the routes for the recipes endpoints
func SetupRecipesRouter(controller *RecipeController, router *gin.RouterGroup) *gin.RouterGroup {
	router.POST("/recipes", controller.CreateRecipeHandler)
	router.GET("/recipes", controller.GetRecipesByFilterHandler)
	router.GET("/recipes/count", controller.CountRecipeByFilterHandler)
	router.GET("/recipes/:id", controller.GetRecipeByIdHandler)
	router.PATCH("/recipes/:id", controller.EditRecipeHandler)
	router.DELETE("/recipes/:id", controller.DeleteRecipeHandler)
	return router
}
