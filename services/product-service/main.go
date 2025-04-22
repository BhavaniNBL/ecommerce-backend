// package main

// import (
// 	"log"
// 	"net"

// 	config "github.com/BhavaniNBL/ecommerce-backend/config/redis"
// 	"github.com/BhavaniNBL/ecommerce-backend/proto/productpb"
// 	"github.com/BhavaniNBL/ecommerce-backend/services/product-service/handler"
// 	"github.com/BhavaniNBL/ecommerce-backend/services/product-service/model"
// 	"github.com/BhavaniNBL/ecommerce-backend/services/product-service/repository"
// 	"github.com/BhavaniNBL/ecommerce-backend/services/product-service/service"
// 	"github.com/gin-gonic/gin"
// 	"google.golang.org/grpc"
// 	"gorm.io/driver/postgres"
// 	"gorm.io/gorm"
// )

// func main() {
// 	dsn := "host=localhost user=postgres password=postgres dbname=productdb port=5432 sslmode=disable"
// 	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		log.Fatal("failed to connect to DB:", err)
// 	}
// 	db.AutoMigrate(&model.Product{})
// 	config.InitRedis()

// 	repo := repository.NewProductRepo(db)
// 	grpcSvc := service.NewProductService(repo)
// 	httpHandler := handler.NewProductHandler(grpcSvc)

// 	// gRPC server
// 	go func() {
// 		lis, err := net.Listen("tcp", ":50052")
// 		if err != nil {
// 			log.Fatal("failed to listen:", err)
// 		}
// 		grpcServer := grpc.NewServer()
// 		productpb.RegisterProductServiceServer(grpcServer, grpcSvc)
// 		log.Println("gRPC server listening on :50052")
// 		if err := grpcServer.Serve(lis); err != nil {
// 			log.Fatal("gRPC server error:", err)
// 		}
// 	}()

// 	// Gin HTTP server
// 	r := gin.Default()
// 	r.POST("/products", httpHandler.CreateProduct)
// 	r.GET("/products/:id", httpHandler.GetProduct)
// 	r.GET("/products", httpHandler.ListProducts)
// 	r.PUT("/products/:id", httpHandler.UpdateProduct)
// 	r.DELETE("/products/:id", httpHandler.DeleteProduct)
// 	log.Println("HTTP server listening on :8082")
// 	r.Run(":8082")
// }

package main

import (
	"context"
	"fmt"
	"log"
	"net"

	appconfig "github.com/BhavaniNBL/ecommerce-backend/config"
	config "github.com/BhavaniNBL/ecommerce-backend/config/redis"
	"github.com/BhavaniNBL/ecommerce-backend/proto/productpb"
	"github.com/BhavaniNBL/ecommerce-backend/services/product-service/handler"
	"github.com/BhavaniNBL/ecommerce-backend/services/product-service/middleware"
	"github.com/BhavaniNBL/ecommerce-backend/services/product-service/model"
	"github.com/BhavaniNBL/ecommerce-backend/services/product-service/repository"
	"github.com/BhavaniNBL/ecommerce-backend/services/product-service/service"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := appconfig.LoadConfig()

	// ---- PostgreSQL Init ----
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Failed to connect to DB:", err)
	}
	db.Migrator().DropTable(&model.Product{}) // WARNING: Deletes all data
	db.AutoMigrate(&model.Product{})

	//db.AutoMigrate(&model.Product{})

	// ---- Redis Init ----
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
	})
	ctx := context.Background()
	err = rdb.Ping(ctx).Err()
	if err != nil {
		log.Fatalf("❌ Redis connection error: %v", err)
	}
	log.Println("✅ Redis connected successfully")

	// 🔥 Add this to make it available globally
	config.RedisClient = rdb
	config.Ctx = ctx

	// ---- Setup Layers ----
	repo := repository.NewProductRepo(db)
	grpcSvc := service.NewProductService(repo)
	httpHandler := handler.NewProductHandler(grpcSvc)

	// ---- gRPC server ----
	go func() {
		lis, err := net.Listen("tcp", ":50052")
		if err != nil {
			log.Fatal("failed to listen:", err)
		}
		grpcServer := grpc.NewServer()
		productpb.RegisterProductServiceServer(grpcServer, grpcSvc)
		log.Println("🚀 gRPC server listening on :50052")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal("gRPC server error:", err)
		}
	}()

	// ---- HTTP server ----
	r := gin.Default()
	r.Use(middleware.JwtMiddleware())
	r.POST("/products", httpHandler.CreateProduct)
	r.GET("/products/:id", httpHandler.GetProduct)
	r.GET("/products", httpHandler.ListProducts)
	r.PUT("/products/:id", httpHandler.UpdateProduct)
	r.DELETE("/products/:id", httpHandler.DeleteProduct)
	log.Println("🚀 HTTP server listening on :8082")
	r.Run(":8082")
}
