package server

import (
	"github.com/gin-gonic/gin"
	"log"
)

func StartWebServer() {
	r := gin.Default()
	InitRoutes(r)
	err := r.Run()
	if err != nil {
		log.Fatal(err)
	}
}
