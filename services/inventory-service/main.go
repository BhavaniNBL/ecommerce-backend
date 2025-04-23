// package main

// import (
// 	"log"
// 	"net"

// 	cfg "github.com/BhavaniNBL/ecommerce-backend/config"
// 	"github.com/BhavaniNBL/ecommerce-backend/config/db"
// 	dbcfg "github.com/BhavaniNBL/ecommerce-backend/config/db"
// 	"github.com/BhavaniNBL/ecommerce-backend/services/inventory-service/handler"
// 	"github.com/BhavaniNBL/ecommerce-backend/services/inventory-service/repository"
// 	"github.com/BhavaniNBL/ecommerce-backend/services/inventory-service/service"
// 	"github.com/BhavaniNBL/ecommerce-backend/shared/cache"
// 	kafka "github.com/BhavaniNBL/ecommerce-backend/shared/kafkautil"

// 	pb "github.com/BhavaniNBL/ecommerce-backend/proto/inventorypb"
// 	"github.com/gin-gonic/gin"
// 	"google.golang.org/grpc"
// )

// func main() {
// 	// Load configuration
// 	cfg := cfg.LoadConfig()

// 	// Initialize DB
// 	dbcfg.InitDB()

// 	// Initialize Redis
// 	redisClient := cache.InitRedisClient()

// 	// Initialize Kafka Producer
// 	// kafkaProducer := kafka.InitKafkaProducer("localhost:9092", "inventory-events")
// 	kafkaProducer := kafka.InitKafkaProducer(cfg.KafkaBroker, cfg.KafkaTopic)
// 	inventoryService := service.NewInventoryService(kafkaProducer, redisClient, cfg.KafkaTopic)

// 	// Initialize Repository
// 	inventoryRepo := repository.NewInventoryRepo(db.DB)

// 	// Initialize Inventory Service
// 	// inventoryService := service.NewInventoryService(kafkaProducer, redisClient)
// 	inventoryService.SetRepo(inventoryRepo)

// 	// Setup gRPC Server
// 	grpcServer := grpc.NewServer()
// 	pb.RegisterInventoryServiceServer(grpcServer, inventoryService)

// 	// Start gRPC server in a separate goroutine
// 	go func() {
// 		lis, err := net.Listen("tcp", ":50053")
// 		if err != nil {
// 			log.Fatalf("❌ Failed to listen on port 50053: %v", err)
// 		}

// 		log.Println("✅ gRPC Server started on port 50053")
// 		if err := grpcServer.Serve(lis); err != nil {
// 			log.Fatalf("❌ Failed to serve gRPC server: %v", err)
// 		}
// 	}()

// 	// Setup HTTP Server (Gin)
// 	r := gin.Default()

// 	// Register HTTP routes (using Gin for handling HTTP requests)
// 	handler.RegisterRoutes(r, inventoryService)

//		// Start HTTP server
//		if err := r.Run(":8083"); err != nil {
//			log.Fatalf("❌ Failed to start HTTP server: %v", err)
//		}
//	}
package main

import (
	"log"
	"net"

	cfg "github.com/BhavaniNBL/ecommerce-backend/config"
	"github.com/BhavaniNBL/ecommerce-backend/config/db"
	"github.com/BhavaniNBL/ecommerce-backend/services/inventory-service/handler"
	"github.com/BhavaniNBL/ecommerce-backend/services/inventory-service/repository"
	"github.com/BhavaniNBL/ecommerce-backend/services/inventory-service/service"
	"github.com/BhavaniNBL/ecommerce-backend/shared/cache"
	kafka "github.com/BhavaniNBL/ecommerce-backend/shared/kafkautil"

	pb "github.com/BhavaniNBL/ecommerce-backend/proto/inventorypb"
	productpb "github.com/BhavaniNBL/ecommerce-backend/proto/productpb"
	middleware "github.com/BhavaniNBL/ecommerce-backend/services/product-service/middleware/interceptor"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Load configuration
	cfg := cfg.LoadConfig()

	// Initialize DB
	db.InitDB()

	// Initialize Redis
	redisClient := cache.InitRedisClient()

	// Initialize Kafka Producer
	kafkaProducer := kafka.InitKafkaProducer(cfg.KafkaBroker, cfg.KafkaTopic)

	// Setup gRPC connection to Product Service
	productConn, err := grpc.NewClient(cfg.ProductServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("❌ Failed to connect to Product Service: %v", err)
	}
	defer productConn.Close()

	productClient := productpb.NewProductServiceClient(productConn)

	// Initialize Inventory Service with Product Client
	inventoryService := service.NewInventoryService(kafkaProducer, redisClient, cfg.KafkaTopic, productClient)

	// Initialize Repository
	inventoryRepo := repository.NewInventoryRepo(db.DB)

	// Set repository for Inventory Service
	inventoryService.SetRepo(inventoryRepo)

	// Setup gRPC Server
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.AuthInterceptor()),
	)
	pb.RegisterInventoryServiceServer(grpcServer, inventoryService)

	// Start gRPC server in a separate goroutine
	go func() {
		lis, err := net.Listen("tcp", ":50053")
		if err != nil {
			log.Fatalf("❌ Failed to listen on port 50053: %v", err)
		}

		log.Println("✅ gRPC Server started on port 50053")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("❌ Failed to serve gRPC server: %v", err)
		}
	}()

	// Setup HTTP Server (Gin)
	r := gin.Default()

	// Register HTTP routes (using Gin for handling HTTP requests)
	handler.RegisterRoutes(r, inventoryService)

	// Start HTTP server
	if err := r.Run(":8083"); err != nil {
		log.Fatalf("❌ Failed to start HTTP server: %v", err)
	}
}
