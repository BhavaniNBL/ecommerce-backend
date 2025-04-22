package main

import (
	"log"
	"net"

	"github.com/BhavaniNBL/ecommerce-backend/config"
	"github.com/BhavaniNBL/ecommerce-backend/config/db"
	"github.com/BhavaniNBL/ecommerce-backend/services/inventory-service/handler"
	"github.com/BhavaniNBL/ecommerce-backend/services/inventory-service/repository"
	"github.com/BhavaniNBL/ecommerce-backend/services/inventory-service/service"
	"github.com/BhavaniNBL/ecommerce-backend/shared/cache"
	"github.com/BhavaniNBL/ecommerce-backend/shared/kafka"

	pb "github.com/BhavaniNBL/ecommerce-backend/proto/inventorypb"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize DB
	db.InitDB()

	// Initialize Redis
	redisClient := cache.InitRedisClient()

	// Initialize Kafka Producer
	// kafkaProducer := kafka.InitKafkaProducer("localhost:9092", "inventory-events")
	kafkaProducer := kafka.InitKafkaProducer(cfg.KafkaBroker)
    inventoryService := service.NewInventoryService(kafkaProducer, redisClient, cfg.KafkaTopic)


	// Initialize Repository
	inventoryRepo := repository.NewInventoryRepo(db.DB)

	// Initialize Inventory Service
	// inventoryService := service.NewInventoryService(kafkaProducer, redisClient)
	inventoryService.SetRepo(inventoryRepo)

	// Setup gRPC Server
	grpcServer := grpc.NewServer()
	pb.RegisterInventoryServiceServer(grpcServer, inventoryService)

	// Start gRPC server in a separate goroutine
	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("❌ Failed to listen on port 50051: %v", err)
		}

		log.Println("✅ gRPC Server started on port 50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("❌ Failed to serve gRPC server: %v", err)
		}
	}()

	// Setup HTTP Server (Gin)
	r := gin.Default()

	// Register HTTP routes (using Gin for handling HTTP requests)
	handler.RegisterRoutes(r, inventoryService)

	// Start HTTP server
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("❌ Failed to start HTTP server: %v", err)
	}
}
