package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/BhavaniNBL/ecommerce-backend/config"
	"github.com/BhavaniNBL/ecommerce-backend/services/payment-service/handler"
	"github.com/BhavaniNBL/ecommerce-backend/services/payment-service/model"
	"github.com/BhavaniNBL/ecommerce-backend/services/payment-service/repository"
	"github.com/BhavaniNBL/ecommerce-backend/services/payment-service/service"
	kafkaUtil "github.com/BhavaniNBL/ecommerce-backend/shared/KafkaConsumerUtil"
	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.LoadConfig()
	dsn := "host=" + cfg.DBHost + " user=" + cfg.DBUser + " password=" + cfg.DBPassword + " dbname=" + cfg.DBName + " port=" + cfg.DBPort + " sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ DB connection failed:", err)
	}
	db.AutoMigrate(&model.Payment{})

	repo := repository.NewPaymentRepo(db)
	svc := service.NewPaymentService(cfg.KafkaBroker, cfg.KafkaTopicOrder, "payment-events", repo)
	h := handler.NewPaymentHandler(repo)

	// Start REST server
	r := gin.Default()
	r.GET("/payments/:orderID", h.GetByOrderID)
	go r.Run(":8085")

	// Start Kafka consumer
	consumerGroup := "payment-consumer"
	topic := cfg.KafkaTopicOrder

	kafkaUtil.ConsumeTopic(cfg.KafkaBroker, consumerGroup, []string{topic}, func(msg *sarama.ConsumerMessage) {
		svc.ProcessOrderEvent(msg.Value)
	})

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Println("🛑 Payment service shutting down")
}
