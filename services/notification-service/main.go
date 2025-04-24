package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/BhavaniNBL/ecommerce-backend/config"
	"github.com/BhavaniNBL/ecommerce-backend/services/notification-service/handler"
	"github.com/BhavaniNBL/ecommerce-backend/services/notification-service/repository"
	"github.com/BhavaniNBL/ecommerce-backend/services/notification-service/service"
	kafkaUtil "github.com/BhavaniNBL/ecommerce-backend/shared/KafkaConsumerUtil"
	"github.com/BhavaniNBL/ecommerce-backend/shared/middleware"
	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.LoadConfig()

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisHost + ":" + cfg.RedisPort,
		Password: "",
		DB:       0,
	})

	repo := repository.NewNotificationRepo(rdb)
	svc := service.NewNotificationService(repo, cfg.KafkaBroker, "notification-events")
	h := handler.NewNotificationHandler(repo, svc)

	r := gin.Default()
	r.Use(middleware.JWTMiddleware())
	r.GET("/notifications/:orderID", h.GetByOrderID)
	r.POST("/notifications", h.CreateNotification)
	go r.Run(":8086")

	log.Println("📡 Notification Service listening to Kafka events...")
	kafkaUtil.ConsumeTopic(cfg.KafkaBroker, "notification-consumer", []string{"order-events"}, func(msg *sarama.ConsumerMessage) {
		svc.HandleOrderConfirmedEvent(msg.Value)
	})

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
	log.Println("🛑 Notification service shutting down")
}
