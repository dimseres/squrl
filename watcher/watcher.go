package watcher

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"groupsCrawl/cmd"
	"groupsCrawl/config"
	"groupsCrawl/watcher/models"
	"groupsCrawl/watcher/parser"
	"log"
)

var (
	ctx = context.Background()
)

type WatchService struct {
	Connection *config.RDB
}

func (service *WatchService) Start() {
	rdb := service.Connection.Connection
	subscriber := rdb.Subscribe(ctx, "send-link")
	messageBus := make(chan models.ChanBus)
	redisBus := make(chan string)
	vkParser := parser.VkParser{
		Urls: &models.PriorityQueue{},
		Bus:  messageBus,
	}

	go vkParser.StartParser()

	go func() {
		for {
			msg, err := subscriber.ReceiveMessage(service.Connection.Cnt)
			if err != nil {
				log.Fatal(err)
			}
			redisBus <- msg.Payload
			formatMessage(msg)
		}
	}()

	for {
		select {
		case redMsg := <-redisBus:
			vkParser.AddLink(redMsg)
		case msg := <-messageBus:
			fmt.Println("Service:", msg.Service, "message:", msg.Message, "payload:", msg.Payload)
		}
	}
}

func (service *WatchService) GetLinks() {
	db, cnt := service.Connection.Connection, service.Connection.Cnt
	db.Get(cnt, "")
}

func formatMessage(message *redis.Message) {
	fmt.Printf(cmd.PurpleBg+"<<< WATCHER >>>"+cmd.Reset+"\tRecieve new watch link from <%s>\n", message.Channel)
}
