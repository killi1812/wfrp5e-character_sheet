package charactersheet

import (
	"errors"
	"net/http"

	"github.com/killi1812/wfrp5e-character_sheet/app"
	"github.com/killi1812/wfrp5e-character_sheet/util/auth"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type CharacterSheetCtn struct {
	service ICharacterSheetService
	logger  *zap.SugaredLogger
}

func NewCharacterSheetCtn() app.Controller {
	var controller *CharacterSheetCtn
	app.Invoke(func(service ICharacterSheetService, logger *zap.SugaredLogger) {
		controller = &CharacterSheetCtn{
			service: service,
			logger:  logger,
		}
	})
	return controller
}

func (c *CharacterSheetCtn) RegisterEndpoints(api *gin.RouterGroup) {
	group := api.Group("/character-sheets")

	// Protected routes
	group.GET("/", auth.Protect(), c.getAll)
	group.GET("/:uuid", auth.Protect(), c.getOne)
	group.POST("/", auth.Protect(), c.create)
	group.PUT("/:uuid", auth.Protect(), c.update)
	group.DELETE("/:uuid", auth.Protect(), c.delete)
}

func (c *CharacterSheetCtn) getAll(ctx *gin.Context) {
	sheets, err := c.service.ReadAll()
	if err != nil {
		c.logger.Errorf("Failed to fetch character sheets: %v", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, sheets)
}

func (c *CharacterSheetCtn) getOne(ctx *gin.Context) {
	sheetUuid, err := uuid.Parse(ctx.Param("uuid"))
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	sheet, err := c.service.Read(sheetUuid)
	if err != nil {
		if errors.Is(err, ErrSheetNotFound) {
			ctx.AbortWithError(http.StatusNotFound, err)
			return
		}
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, sheet)
}

func (c *CharacterSheetCtn) create(ctx *gin.Context) {
	var sheet CharacterSheet
	if err := ctx.BindJSON(&sheet); err != nil {
		c.logger.Errorf("Failed to bind character sheet JSON: %v", err)
		return
	}

	// Try extracting logged in user ID from token
	if authHeader := ctx.GetHeader("Authorization"); authHeader != "" {
		if _, claims, err := auth.ParseToken(authHeader); err == nil && claims.ID != "" {
			if parsedUserUuid, err := uuid.Parse(claims.ID); err == nil {
				sheet.UserUuid = parsedUserUuid
			}
		}
	}

	created, err := c.service.Create(&sheet)
	if err != nil {
		c.logger.Errorf("Failed to create character sheet: %v", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, created)
}

func (c *CharacterSheetCtn) update(ctx *gin.Context) {
	sheetUuid, err := uuid.Parse(ctx.Param("uuid"))
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var sheet CharacterSheet
	if err := ctx.BindJSON(&sheet); err != nil {
		c.logger.Errorf("Failed to bind update JSON: %v", err)
		return
	}

	updated, err := c.service.Update(sheetUuid, &sheet)
	if err != nil {
		if errors.Is(err, ErrSheetNotFound) {
			ctx.AbortWithError(http.StatusNotFound, err)
			return
		}
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, updated)
}

func (c *CharacterSheetCtn) delete(ctx *gin.Context) {
	sheetUuid, err := uuid.Parse(ctx.Param("uuid"))
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err = c.service.Delete(sheetUuid)
	if err != nil {
		if errors.Is(err, ErrSheetNotFound) {
			ctx.AbortWithError(http.StatusNotFound, err)
			return
		}
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
