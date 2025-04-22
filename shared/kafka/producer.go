package kafka

import (
	"log"

	"github.com/segmentio/kafka-go"
)

var Writer *kafka.Writer

// Initialize the Kafka Producer
func InitKafkaProducer(broker string, topic string) *kafka.Writer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(broker),   // Kafka broker address
		Topic:    topic,               // Topic to publish messages
		Balancer: &kafka.LeastBytes{}, // Balancing strategy for partitioning
	}

	Writer = writer

	log.Println("✅ Kafka Producer initialized")
	return writer
}

// PublishMessage sends a message to Kafka topic
func PublishMessage(writer *kafka.Writer, key string, value string) error {
	message := kafka.Message{
		Key:   []byte(key),
		Value: []byte(value),
	}

	// Send message to Kafka topic
	err := writer.WriteMessages(nil, message)
	if err != nil {
		log.Printf("❌ Error publishing message to Kafka: %v", err)
		return err
	}

	log.Println("✅ Kafka message published")
	return nil
}
