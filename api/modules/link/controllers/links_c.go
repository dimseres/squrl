package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"groupsCrawl/api/modules/link/models"
	"groupsCrawl/api/modules/link/services"
	"log"
)

func UpdateLink(ctx *gin.Context) {
	var UriParams = struct {
		Id string `uri:"id" binding:"required,uuid"`
	}{}
	if err := ctx.ShouldBindUri(&UriParams); err != nil {
		ctx.JSON(400, gin.H{"msg": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"id": UriParams.Id})
}

func AddNewLink(ctx *gin.Context) {
	var body models.AddLinkRequest
	jsonBytes, err := ctx.GetRawData()
	if err != nil {
		log.Fatalln(err)
	}
	err = json.Unmarshal(jsonBytes, &body)
	if err != nil {
		log.Fatalln(err)
	}
	srvc := services.NewLinkService(ctx)
	ctx.JSON(srvc.AddNewLink(body))
}
