// package config

// import (
// 	"log"
// 	"os"

// 	"github.com/joho/godotenv"
// )

// func LoadConfig() {
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatalf("Error loading .env file")
// 	}

// 	log.Println("Configuration loaded successfully")
// }

// func GetEnv(key, fallback string) string {
// 	value, exists := os.LookupEnv(key)
// 	if !exists {
// 		return fallback
// 	}
// 	return value
// }

package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	RedisHost  string
	RedisPort  string
	KafkaHost  string
	KafkaPort  string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("❌ Error loading .env file")
	}

	log.Println("✅ Configuration loaded successfully")

	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "ecommerce_db"),
		RedisHost:  getEnv("REDIS_HOST", "localhost"),
		RedisPort:  getEnv("REDIS_PORT", "6379"),
		KafkaHost:  getEnv("KAFKA_HOST", "localhost"),
		KafkaPort:  getEnv("KAFKA_PORT", "9092"),
	}
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}
