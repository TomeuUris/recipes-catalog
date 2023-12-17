package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/TomeuUris/recipes-catalog/api/v1/payload"
	"github.com/TomeuUris/recipes-catalog/pkg/entity"
	"github.com/TomeuUris/recipes-catalog/pkg/repo/ingredient"
	"github.com/gin-gonic/gin"
)

type IngredientController struct {
	repo ingredient.Repo
}

func NewController(repo ingredient.Repo) *IngredientController {
	return &IngredientController{repo: repo}
}

// @Summary Get ingredient by ID
// @Description get User by its ID
// @Tags Ingredients
// @Produce  json
// @Param   id     path    int     true        "Ingredient ID"
// @Success 200 {object} entity.Ingredient
// @Router /ingredients/{id} [get]
func (c *IngredientController) GetIngredientByIdHandler(ctx *gin.Context) {
	// Get the ingredient ID from the URL parameter
	ingredientIDStr := ctx.Params.ByName("id")
	ingredientID, err := strconv.Atoi(ingredientIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ingredient ID"})
		return
	}

	// Get the ingredient from the database
	// dbctx := ctx.Request.Context()
	ingredient, err := c.repo.FindByID(ctx, ingredientID)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the ingredients as a response
	ctx.JSON(http.StatusOK, ingredient)
}

// @Summary Create ingredient
// @Description Create a new Ingredient
// @Tags Ingredients
// @Accept  json
// @Produce  json
// @Param   ingredient     body    payload.Ingredient     true        "Ingredient info"
// @Success 200 {object} entity.Ingredient
// @Router /ingredients [post]
func (c *IngredientController) CreateIngredientHandler(ctx *gin.Context) {
	// Parse the request payload
	var ingredientPayload payload.Ingredient
	if err := ctx.ShouldBindJSON(&ingredientPayload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ingredient := ingredientPayload.Entity()
	// Logic to create the ingredient in the database
	err := c.repo.Add(ctx, ingredient)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the created ingredient as a response
	ctx.JSON(http.StatusCreated, ingredient)
}

func SetupRouter(controller *IngredientController, router *gin.Engine) *gin.Engine {
	router.GET("/ingredients/:id", controller.GetIngredientByIdHandler)
	router.POST("/ingredients", controller.CreateIngredientHandler)
	return router
}
