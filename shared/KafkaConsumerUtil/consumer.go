package kafkaConsumerUtil

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/IBM/sarama"
)

type consumerGroupHandler struct {
	handlerFunc func(*sarama.ConsumerMessage)
}

func (h consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (h consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }
func (h consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		h.handlerFunc(message)
		session.MarkMessage(message, "")
	}
	return nil
}

func ConsumeTopic(broker, groupID string, topics []string, handlerFunc func(*sarama.ConsumerMessage)) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Version = sarama.V2_1_0_0

	consumer, err := sarama.NewConsumerGroup(strings.Split(broker, ","), groupID, config)
	if err != nil {
		log.Fatalf("❌ Failed to create consumer group: %v", err)
	}

	ctx := context.Background()
	go func() {
		for {
			if err := consumer.Consume(ctx, topics, consumerGroupHandler{handlerFunc}); err != nil {
				log.Printf("❌ Error consuming messages: %v", err)
			}
			time.Sleep(1 * time.Second)
		}
	}()
}
