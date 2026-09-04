package app

import (
	"github.com/gin-gonic/gin"
)

type Controller interface {
	RegisterEndpoints(router *gin.RouterGroup)
}

var apis []Controller

// RegisterApi registers a controller to a router
func RegisterApi(newApi func() Controller) {
	apis = append(apis, newApi())
}
