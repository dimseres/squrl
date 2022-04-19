package services

import (
	"context"
	"github.com/gin-gonic/gin"
	"groupsCrawl/api/modules/link/models"
	"groupsCrawl/config"
	"net/http"
)

var (
	rdb            = config.Redis()
	ctx            = context.Background()
	CreatedService *LinkService
)

type LinkService struct {
	name    string
	context *gin.Context
}

func (service *LinkService) ServiceName() string {
	return service.name
}

func (service *LinkService) AddNewLink(body models.AddLinkRequest) (int, gin.H) {
	rdb = config.Redis()
	rdb.Connection.Publish(rdb.Cnt, "send-link", body.Url)
	return http.StatusOK, gin.H{
		"message": body.Url,
		"ctx":     service.context.Request.Host,
	}
}

func NewLinkService(ctx *gin.Context) *LinkService {
	newService := LinkService{
		name:    "LinkService",
		context: ctx,
	}
	CreatedService = &newService
	return &newService
}

func InitLinkService(ctx *gin.Context) LinkService {
	return LinkService{
		name:    "LinkService",
		context: ctx,
	}
}
