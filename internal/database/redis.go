package database

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"

	"github.com/testgithubanjali/url-shortener/internal/config"
)

var (
	RedisClient *redis.Client
	Ctx         = context.Background()
)

func ConnectRedis() {
	addr := fmt.Sprintf("%s:%s",
		config.AppConfig.RedisHost,
		config.AppConfig.RedisPort,
	)

	log.Printf("Connecting to Redis at %s", addr)

	RedisClient = redis.NewClient(&redis.Options{
		Addr: addr,
	})

	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		log.Fatal("❌ Failed to connect Redis:", err)
	}

	log.Println("Redis connected successfully")
}
