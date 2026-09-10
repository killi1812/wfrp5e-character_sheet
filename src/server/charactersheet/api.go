package charactersheet

import (
	"net/http"

	"github.com/killi1812/wfrp5e-character_sheet/app"
	"github.com/killi1812/wfrp5e-character_sheet/util/auth"
	"github.com/killi1812/wfrp5e-character_sheet/util/ginutil"

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
	group := api.Group("/character-sheets", auth.Protect())

	// Protected routes
	group.GET("/", c.getAll)
	group.GET("/:uuid", c.getOne)
	group.POST("/", c.create)
	group.PUT("/:uuid", c.update)
	group.DELETE("/:uuid", c.delete)
}

func getUserInfo(ctx *gin.Context) (uuid.UUID, string) {
	if authHeader := ctx.GetHeader("Authorization"); authHeader != "" {
		if _, claims, err := auth.ParseToken(authHeader); err == nil && claims.ID != "" {
			if parsedUserUuid, err := uuid.Parse(claims.ID); err == nil {
				return parsedUserUuid, claims.Role
			}
		}
	}
	return uuid.Nil, ""
}

func (c *CharacterSheetCtn) getAll(ctx *gin.Context) {
	userUuid, role := getUserInfo(ctx)

	var sheets []CharacterSheet
	var err error
	if role == "admin" || userUuid == uuid.Nil {
		sheets, err = c.service.ReadAll()
	} else {
		sheets, err = c.service.ReadByUser(userUuid)
	}

	if err != nil {
		c.logger.Errorf("Failed to fetch character sheets: %v", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, sheets)
}

func (c *CharacterSheetCtn) hasAccess(ctx *gin.Context, sheetOwner uuid.UUID) bool {
	userUuid, role := getUserInfo(ctx)
	if role != "admin" && sheetOwner != uuid.Nil && userUuid != uuid.Nil && sheetOwner != userUuid {
		ctx.AbortWithStatus(http.StatusForbidden)
		return false
	}
	return true
}

func (c *CharacterSheetCtn) getOne(ctx *gin.Context) {
	sheetUuid, ok := ginutil.ParseUUIDParam(ctx, "uuid")
	if !ok {
		return
	}

	sheet, err := c.service.Read(sheetUuid)
	if ginutil.HandleServiceError(ctx, err, ErrSheetNotFound) {
		return
	}

	if !c.hasAccess(ctx, sheet.UserUuid) {
		return
	}

	ctx.JSON(http.StatusOK, sheet)
}

func (c *CharacterSheetCtn) create(ctx *gin.Context) {
	var sheet CharacterSheet
	if err := ctx.BindJSON(&sheet); err != nil {
		c.logger.Errorf("Failed to bind character sheet JSON: %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	userUuid, _ := getUserInfo(ctx)
	if userUuid != uuid.Nil {
		sheet.UserUuid = userUuid
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
	sheetUuid, ok := ginutil.ParseUUIDParam(ctx, "uuid")
	if !ok {
		return
	}

	existing, err := c.service.Read(sheetUuid)
	if ginutil.HandleServiceError(ctx, err, ErrSheetNotFound) {
		return
	}

	if !c.hasAccess(ctx, existing.UserUuid) {
		return
	}

	var sheet CharacterSheet
	if err := ctx.BindJSON(&sheet); err != nil {
		c.logger.Errorf("Failed to bind update JSON: %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if sheet.UserUuid == uuid.Nil {
		sheet.UserUuid = existing.UserUuid
	}

	updated, err := c.service.Update(sheetUuid, &sheet)
	if err != nil {
		c.logger.Errorf("Failed to update character sheet: %v", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, updated)
}

func (c *CharacterSheetCtn) delete(ctx *gin.Context) {
	sheetUuid, ok := ginutil.ParseUUIDParam(ctx, "uuid")
	if !ok {
		return
	}

	existing, err := c.service.Read(sheetUuid)
	if ginutil.HandleServiceError(ctx, err, ErrSheetNotFound) {
		return
	}

	if !c.hasAccess(ctx, existing.UserUuid) {
		return
	}

	err = c.service.Delete(sheetUuid)
	if err != nil {
		c.logger.Errorf("Failed to delete character sheet: %v", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
