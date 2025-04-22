package cache

import (
	"context"
	"log"

	"github.com/BhavaniNBL/ecommerce-backend/config"
	"github.com/redis/go-redis/v9"
)

// Global Redis client instance
var Client *redis.Client

// Initialize Redis Client
func InitRedisClient() *redis.Client {
	cfg := config.LoadConfig()

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisHost + ":" + cfg.RedisPort, // Redis connection info
		Password: "",                                  // No password
		DB:       0,                                   // Default DB
	})

	Client = client

	// Test Redis connection
	_, err := Client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("❌ Error connecting to Redis: %v", err)
	}

	log.Println("✅ Redis connection established")
	return client
}
