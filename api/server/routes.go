package server

import (
	"github.com/gin-gonic/gin"
	"groupsCrawl/api/modules/link/routes"
)

func InitRoutes(router *gin.Engine) {
	apiV1 := router.Group("/api/v1")
	routes.AddHandlers(apiV1)
}
