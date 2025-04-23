package kafkautil

import (
	"log"
	"time"

	"github.com/IBM/sarama"
)

var Producer sarama.SyncProducer

const (
	maxInitRetries    = 5
	initialBackoff    = 2 * time.Second
	maxPublishRetries = 3
)

// InitKafkaProducer initializes the Kafka producer using sarama
func InitKafkaProducer(broker, topic string) sarama.SyncProducer {
	var producer sarama.SyncProducer
	var err error

	for attempt := 1; attempt <= maxInitRetries; attempt++ {
		config := sarama.NewConfig()
		config.Producer.Return.Successes = true
		config.Producer.RequiredAcks = sarama.WaitForAll

		producer, err = sarama.NewSyncProducer([]string{broker}, config)
		if err == nil {
			log.Println("✅ Kafka Producer initialized successfully")
			Producer = producer
			return producer
		}

		log.Printf("⚠️ Kafka initialization failed (attempt %d/%d): %v", attempt, maxInitRetries, err)
		time.Sleep(initialBackoff * time.Duration(attempt))
	}

	log.Fatal("❌ Failed to initialize Kafka producer after retries")
	return nil
}

// PublishMessage sends a message to a Kafka topic using sarama
func PublishMessage(producer sarama.SyncProducer, topic, key, value string) error {
	message := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(value),
	}

	var err error
	for attempt := 1; attempt <= maxPublishRetries; attempt++ {
		_, _, err = producer.SendMessage(message)
		if err == nil {
			log.Println("✅ Kafka message published")
			return nil
		}

		log.Printf("⚠️ Kafka publish failed (attempt %d/%d): %v", attempt, maxPublishRetries, err)
		time.Sleep(500 * time.Millisecond * time.Duration(attempt))
	}

	log.Printf("❌ Kafka publish failed after %d attempts: %v", maxPublishRetries, err)
	return err
}

// SafePublish is a safe wrapper around PublishMessage with auto reinitialization
func SafePublish(broker, topic, key, value string) {
	if Producer == nil {
		log.Println("⚠️ Kafka producer not initialized, retrying...")
		Producer = InitKafkaProducer(broker, topic)
	}

	// Check if the topic exists before publishing (you can create it externally)
	if !TopicExists(broker, topic) {
		log.Printf("❌ Topic %s does not exist. Please create it manually or handle the creation externally.", topic)
	}

	// Publish the message to Kafka
	err := PublishMessage(Producer, topic, key, value)
	if err != nil {
		// Optionally store in fallback queue / DB
		log.Printf("❌ Message lost or queued for retry: %v", err)
	}
}

// TopicExists checks if a Kafka topic exists by attempting to read metadata
func TopicExists(broker, topic string) bool {
	// Connect to Kafka and retrieve metadata
	client, err := sarama.NewClient([]string{broker}, nil)
	if err != nil {
		log.Printf("❌ Error connecting to Kafka: %v", err)
		return false
	}
	defer client.Close()

	// Check if the topic exists in Kafka
	topics, err := client.Topics()
	if err != nil {
		log.Printf("❌ Error fetching topics: %v", err)
		return false
	}

	for _, t := range topics {
		if t == topic {
			log.Printf("✅ Topic %s exists.", topic)
			return true
		}
	}

	log.Printf("❌ Topic %s does not exist.", topic)
	return false
}
