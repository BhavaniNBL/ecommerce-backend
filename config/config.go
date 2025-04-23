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
	DBHost             string
	DBPort             string
	DBUser             string
	DBPassword         string
	DBName             string
	RedisHost          string
	RedisPort          string
	KafkaHost          string
	KafkaPort          string
	KafkaBroker        string
	KafkaTopic         string
	KafkaTopicOrder    string
	ProductServiceAddr string
	InventoryServiceAddr string

	// KafkaTopicPayment string
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
		// KafkaHost:  getEnv("KAFKA_HOST", "localhost"),
		// KafkaPort:  getEnv("KAFKA_PORT", "9092"),
		KafkaHost:          getEnv("KAFKA_HOST", "localhost"),
		KafkaPort:          getEnv("KAFKA_PORT", "9092"),
		KafkaBroker:        getEnv("KAFKA_BROKER", "localhost:9092"),
		KafkaTopic:         getEnv("KAFKA_TOPIC_INVENTORY", "inventory-events"),
		KafkaTopicOrder:    getEnv("KAFKA_TOPIC_ORDER", "order-events"),
		ProductServiceAddr: getEnv("PRODUCT_SERVICE_ADDR", "product-service:50052"),
		InventoryServiceAddr: getEnv("INVENTORY_SERVICE_ADDR", "inventory-service:50053"),

		// KafkaTopicPayment: getEnv("KAFKA_TOPIC_PAYMENT", "payment-events"),
	}
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}
