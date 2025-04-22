package cache

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var (
	Rdb *redis.Client
	Ctx = context.Background()
)

func InitRedis() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	_, err := Rdb.Ping(Ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("failed to connect to Redis: %v", err))
	}
}
