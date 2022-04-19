package config

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"os"
)

var ctx = context.Background()

func InitApp() {
	redisURI := os.Getenv("REDIS_URI")
	redisPwd := os.Getenv("REDIS_PASSWORD")

	connection := redis.NewClient(&redis.Options{
		Addr:     redisURI,
		Password: redisPwd,
		DB:       0,
	})

	rdb = RDB{
		Connection: connection,
		Cnt:        ctx,
	}

	fmt.Println("========== INIT COMPLETE ===========")
}

type RDB struct {
	Connection *redis.Client
	Cnt        context.Context
}

var rdb RDB

func Redis() *RDB {
	return &rdb
}
