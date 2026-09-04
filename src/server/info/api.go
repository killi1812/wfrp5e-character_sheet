package info

import (
	"net/http"

	"github.com/killi1812/wfrp5e-character_sheet/app"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type InfoCtn struct {
	logger *zap.SugaredLogger
}

func NewInfoCnt() app.Controller {
	var controller *InfoCtn
	app.Invoke(func(logger *zap.SugaredLogger) {
		controller = &InfoCtn{
			logger: logger,
		}
	})
	return controller
}

func (cnt *InfoCtn) RegisterEndpoints(router *gin.RouterGroup) {
	router.GET("/info", cnt.getServerInfo)
}

// getServerInfo godoc
//
//	@Summary		Get server info
//	@Description	return information about the server build, version, etc ...
//	@Tags			info
//	@Produce		json
//	@Success		200	{object}	ServerInfoDto	"Information about server"
//	@Router			/info [get]
func (ctn *InfoCtn) getServerInfo(c *gin.Context) {
	serverInfo := ServerInfoDto{
		Build:          app.Build,
		Version:        app.Version,
		CommitHash:     app.CommitHash,
		BuildTimestamp: app.BuildTimestamp,
	}
	c.AbortWithStatusJSON(http.StatusOK, serverInfo)
}
