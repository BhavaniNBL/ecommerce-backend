package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/BhavaniNBL/ecommerce-backend/config"
	inventorypb "github.com/BhavaniNBL/ecommerce-backend/proto/inventorypb"
	"github.com/BhavaniNBL/ecommerce-backend/services/order-service/handler"
	"github.com/BhavaniNBL/ecommerce-backend/services/order-service/model"
	"github.com/BhavaniNBL/ecommerce-backend/services/order-service/repository"
	"github.com/BhavaniNBL/ecommerce-backend/services/order-service/service"
	"github.com/BhavaniNBL/ecommerce-backend/shared/middleware"
)

func main() {
	cfg := config.LoadConfig()
	dsn := "host=" + cfg.DBHost + " user=" + cfg.DBUser + " password=" + cfg.DBPassword + " dbname=" + cfg.DBName + " port=" + cfg.DBPort + " sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ DB connection failed:", err)
	}
	db.AutoMigrate(&model.Order{})

	// conn, err := grpc.Dail(cfg.InventoryServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial(cfg.InventoryServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatal("❌ gRPC connection failed:", err)
	}

	repo := repository.NewOrderRepository(db)
	invClient := inventorypb.NewInventoryServiceClient(conn)
	svc := service.NewOrderService(repo, invClient, cfg.KafkaBroker, cfg.KafkaTopicOrder)
	handler := handler.NewOrderHandler(svc)

	r := gin.Default()
	r.Use(middleware.JWTMiddleware())
	r.POST("/orders", handler.CreateOrder)
	r.GET("/orders", handler.ListOrders)
	r.GET("/orders/:id", handler.GetOrder)
	r.PUT("/orders/:id/status", handler.UpdateOrderStatus)

	log.Println("🚀 Order service listening on :8081")
	r.Run(":8081")
}
