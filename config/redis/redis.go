package config

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
	Ctx         = context.Background()
)

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "redis:6379", // or use env/config
		Password: "",
		DB:       0,
	})

	// _, err := RedisClient.Ping(Ctx).Result()
	// if err != nil {
	// 	log.Fatalf("Could not connect to Redis: %v", err)
	// }
	// log.Println("Connected to Redis")

	if err := RedisClient.Ping(Ctx).Err(); err != nil {
		panic("❌ Failed to connect to Redis: " + err.Error())
	}
	log.Println("Connected to Redis")
}
