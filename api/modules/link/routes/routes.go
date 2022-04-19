package routes

import (
	"github.com/gin-gonic/gin"
	"groupsCrawl/api/modules/link/controllers"
)

func AddHandlers(route *gin.RouterGroup) {
	linkRoute := route.Group("link")
	linkRoute.Use(func(context *gin.Context) {
		controllers.InitServices(context)
	})
	linkRoute.POST("/:id/update", controllers.UpdateLink)
	linkRoute.POST("/add", controllers.AddNewLink)
	linkRoute.GET("", func(context *gin.Context) {
		//time.Sleep(time.Second * 50)
		context.JSON(200, gin.H{
			"message": "PK",
		})
	})
}
