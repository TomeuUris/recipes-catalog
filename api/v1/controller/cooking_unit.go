package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/TomeuUris/recipes-catalog/api/v1/payload"
	"github.com/TomeuUris/recipes-catalog/api/v1/view"
	"github.com/TomeuUris/recipes-catalog/pkg/entity"
	"github.com/TomeuUris/recipes-catalog/pkg/repo/cooking_unit"
	"github.com/gin-gonic/gin"
)

type CookingUnitController struct {
	repo cooking_unit.Repo
}

func NewCookingUnitController(repo cooking_unit.Repo) *CookingUnitController {
	return &CookingUnitController{repo: repo}
}

// @Summary Get cooking unit by filter
// @Description Retrieves a list of cooking units filtered by the given parameters
// @Tags Cooking Units
// @Produce  json
// @Param   filter     query    cooking_unit.FindFilter     true        "Filter parameters"
// @Success 200 {object} view.CookingUnit
// @Router /cooking-units [get]
func (c *CookingUnitController) GetCookingUnitByFilterHandler(ctx *gin.Context) {
	// Parse the filter from the query parameters
	var filter cooking_unit.FindFilter
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find the ingredients in the database
	cooking_units, err := c.repo.FindByFilter(ctx, &filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert the ingredients to a view
	cookingUnitViews := make([]*view.CookingUnit, len(cooking_units))
	for i, cooking_unit := range cooking_units {
		cookingUnitViews[i] = &view.CookingUnit{}
		cookingUnitViews[i].FromEntity(cooking_unit)
	}

	// Return the ingredients as a response
	ctx.JSON(http.StatusOK, cookingUnitViews)
}

// @Summary Count cooking units by filter
// @Description Retrieves the number of cooking units filtered by the given parameters
// @Tags Cooking Units
// @Produce  json
// @Param   filter     query    cooking_unit.FindFilter     true        "Filter parameters"
// @Success 200 {object} int
// @Router /cooking-units/count [get]
func (c *CookingUnitController) CountCookingUnitByFilterHandler(ctx *gin.Context) {
	// Parse the filter from the query parameters
	var filter cooking_unit.FindFilter
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Count the ingredients in the database
	count, err := c.repo.CountByFilter(&filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the ingredients as a response
	ctx.JSON(http.StatusOK, count)
}

// @Summary Get cooking unit by ID
// @Description Retrieves an Cooking Unit by its ID
// @Tags Cooking Units
// @Produce  json
// @Param   id     path    int     true        "Cooking Unit ID"
// @Success 200 {object} view.CookingUnit
// @Router /cooking-units/{id} [get]
func (c *CookingUnitController) GetCookingUnitByIdHandler(ctx *gin.Context) {
	// Get the ingredient ID from the URL parameter
	cookingUnitIDStr := ctx.Params.ByName("id")
	cookingUnitID, err := strconv.Atoi(cookingUnitIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ingredient ID"})
		return
	}

	// Get the ingredient from the database
	cookingUnit, err := c.repo.FindByID(ctx, cookingUnitID)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert the cooking unit to a view
	cookingUnitView := &view.CookingUnit{}
	cookingUnitView.FromEntity(cookingUnit)

	// Return the ingredient as a response
	ctx.JSON(http.StatusOK, cookingUnitView)
}

// @Summary Create cooking unit
// @Description Create a new Cooking Unit
// @Tags Cooking Units
// @Accept  json
// @Produce  json
// @Param   cookingUnit     body    payload.CookingUnit     true        "Cooking unit info"
// @Success 200 {object} view.CookingUnit
// @Router /cooking-units [post]
func (c *CookingUnitController) CreateCookingUnitHandler(ctx *gin.Context) {
	// Parse the request payload
	var cookingUnitPayload payload.CookingUnit
	if err := ctx.ShouldBindJSON(&cookingUnitPayload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert the payload to an entity
	cookingUnit := cookingUnitPayload.ToEntity()

	// Create the cooking unit in the database
	err := c.repo.Add(ctx, cookingUnit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert the cooking unit to a view
	cookingUnitView := &view.CookingUnit{}
	cookingUnitView.FromEntity(cookingUnit)

	// Return the created cooking unit as a response
	ctx.JSON(http.StatusCreated, cookingUnitView)
}

// @Summary Edit cooking unit
// @Description Edits an existing cooking unit
// @Tags Cooking Units
// @Accept  json
// @Produce  json
// @Param   id     path    int64               true        "Cooking unit ID"
// @Param   cookingUnit     body    payload.CookingUnit     true        "Cooking unit info"
// @Success 200 {object} view.CookingUnit
// @Router /cooking-units/{id} [patch]
func (c *CookingUnitController) EditCookingUnitHandler(ctx *gin.Context) {
	// Get the cooking unit ID from the URL parameter
	cookingUnitIDStr := ctx.Params.ByName("id")
	cookingUnitID, err := strconv.ParseInt(cookingUnitIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cooking unit ID"})
		return
	}
	// Parse the request payload
	var cookingUnitPayload payload.CookingUnit
	if err := ctx.ShouldBindJSON(&cookingUnitPayload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find the cooking unit in the database
	targetCookingUnit, err := c.repo.FindByID(ctx, int(cookingUnitID))
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Apply the changes to the cooking unit
	cookingUnitPayload.ApplyTo(targetCookingUnit)

	// Update the cooking unit in the database
	err = c.repo.Edit(ctx, targetCookingUnit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert the cooking unit to a view
	cookingUnitView := &view.CookingUnit{}
	cookingUnitView.FromEntity(targetCookingUnit)

	// Return the updated cooking unit as a response
	ctx.JSON(http.StatusOK, cookingUnitView)
}

// @Summary Delete cooking unit
// @Description Deletes an existing Cooking Unit
// @Tags Cooking Units
// @Accept  json
// @Produce  json
// @Param   id     path    int64               true        "Cooking unit ID"
// @Success 204 "No Content"
// @Router /cooking-units/{id} [delete]
func (c *CookingUnitController) DeleteCookingUnitHandler(ctx *gin.Context) {
	// Get the ingredient ID from the URL parameter
	cookingUnitIDStr := ctx.Params.ByName("id")
	cookingUnitID, err := strconv.ParseInt(cookingUnitIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ingredient ID"})
		return
	}

	cookingUnitEntity := &entity.CookingUnit{ID: cookingUnitID}
	// Logic to create the ingredient in the database
	err = c.repo.Delete(ctx, cookingUnitEntity)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the updated ingredient as a response
	ctx.JSON(http.StatusNoContent, nil)
}

func SetupCookingUnitsRouter(controller *CookingUnitController, router *gin.RouterGroup) *gin.RouterGroup {
	router.GET("/cooking-units", controller.GetCookingUnitByFilterHandler)
	router.POST("/cooking-units", controller.CreateCookingUnitHandler)
	router.GET("/cooking-units/count", controller.CountCookingUnitByFilterHandler)
	router.GET("/cooking-units/:id", controller.GetCookingUnitByIdHandler)
	router.PATCH("/cooking-units/:id", controller.EditCookingUnitHandler)
	router.DELETE("/cooking-units/:id", controller.DeleteCookingUnitHandler)
	return router
}
