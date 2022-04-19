package controllers

import (
	"github.com/gin-gonic/gin"
	"groupsCrawl/api/modules/link/services"
)

var linkService *services.LinkService

func InitServices(ctx *gin.Context) {
	linkService = services.NewLinkService(ctx)
}
