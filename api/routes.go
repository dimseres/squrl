package api

import (
	"github.com/gin-gonic/gin"
	"groupsCrawl/api/modules/link/routes"
	"groupsCrawl/ws"
)

func InitRoutes(router *gin.Engine) {
	router.GET("/ws", func(context *gin.Context) {
		ws.Echo(context.Writer, context.Request)
	})
	apiV1 := router.Group("/api/v1")

	routes.AddHandlers(apiV1)
}
