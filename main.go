package main

import (
	"github.com/joho/godotenv"
	"groupsCrawl/config"
	"groupsCrawl/config/server"
	"groupsCrawl/watcher"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln(err)
	}
	config.InitApp()

	rdb := config.Redis()
	watchService := watcher.WatchService{
		Connection: rdb,
	}

	rdb.Connection.Subscribe(rdb.Cnt, "send-link")

	go watchService.Start()
	server.StartWebServer()
}
