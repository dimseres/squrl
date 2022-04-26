package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"groupsCrawl/api"
	"log"
	"net/http"
	"time"
)

func StartWebServer() {
	r := gin.Default()

	r.Static("/storage", "./storage")

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api.InitRoutes(r)
	err := http.ListenAndServe("localhost:8000", r)
	if err != nil {
		log.Fatal(err)
	}
}
