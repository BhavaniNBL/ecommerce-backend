// package service

// import (
// 	"context"
// 	"encoding/json"

// 	"github.com/BhavaniNBL/ecommerce-backend/services/inventory-service/model"
// 	"github.com/BhavaniNBL/ecommerce-backend/services/inventory-service/repository"

// 	// kafkautil "github.com/BhavaniNBL/ecommerce-backend/shared/kafkautil"
// 	kafkautil "github.com/BhavaniNBL/ecommerce-backend/shared/kafkautil"
// 	segmentioKafka "github.com/segmentio/kafka-go"

// 	pb "github.com/BhavaniNBL/ecommerce-backend/proto/inventorypb"
// 	"github.com/redis/go-redis/v9"
// )

// type InventoryService struct {
// 	pb.UnimplementedInventoryServiceServer
// 	repo        repository.InventoryRepo
// 	kafkaWriter *segmentioKafka.Writer
// 	redis       *redis.Client
// 	topicName   string
// }

// func NewInventoryService(kafkaProducer *segmentioKafka.Writer, redisClient *redis.Client, topicName string) *InventoryService {
// 	return &InventoryService{
// 		kafkaWriter: kafkaProducer,
// 		redis:       redisClient,
// 		topicName:   topicName,
// 	}
// }

// func (s *InventoryService) SetRepo(repo repository.InventoryRepo) {
// 	s.repo = repo
// }

// func (s *InventoryService) GetInventory(ctx context.Context, req *pb.GetInventoryRequest) (*pb.InventoryResponse, error) {
// 	// Check Redis cache
// 	val, err := s.redis.Get(ctx, req.ProductId).Result()
// 	if err == nil {
// 		var cached model.Inventory
// 		if err := json.Unmarshal([]byte(val), &cached); err == nil {
// 			return &pb.InventoryResponse{
// 				ProductId:         cached.ProductID,
// 				Quantity:          cached.Quantity,
// 				WarehouseLocation: cached.WarehouseLocation,
// 			}, nil
// 		}
// 	}

// 	inv, err := s.repo.GetInventory(req.ProductId)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Cache in Redis
// 	data, _ := json.Marshal(inv)
// 	s.redis.Set(ctx, req.ProductId, data, 0)

// 	return &pb.InventoryResponse{
// 		ProductId:         inv.ProductID,
// 		Quantity:          inv.Quantity,
// 		WarehouseLocation: inv.WarehouseLocation,
// 	}, nil
// }

// func (s *InventoryService) UpdateInventory(ctx context.Context, req *pb.UpdateInventoryRequest) (*pb.InventoryResponse, error) {
// 	inv, err := s.repo.UpdateInventory(req.ProductId, req.QuantityChange)
// 	if err != nil {
// 		return nil, err
// 	}

// 	eventType := "InventoryReserved"
// 	if inv.Quantity < 0 {
// 		eventType = "OutOfStock"
// 	}

// 	eventPayload := map[string]interface{}{
// 		"type":       eventType,
// 		"product_id": req.ProductId,
// 		"quantity":   inv.Quantity,
// 	}
// 	data, _ := json.Marshal(eventPayload)
// 	//kafkautil.PublishMessage(s.kafkaWriter, "inventory-events", req.ProductId, string(data))
// 	kafkautil.SafePublish("inventory-events", req.ProductId, string(data))

// 	s.redis.Del(ctx, req.ProductId) // invalidate cache

// 	return &pb.InventoryResponse{
// 		ProductId:         inv.ProductID,
// 		Quantity:          inv.Quantity,
// 		WarehouseLocation: inv.WarehouseLocation,
// 	}, nil
// }

package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/BhavaniNBL/ecommerce-backend/services/inventory-service/model"
	"github.com/BhavaniNBL/ecommerce-backend/services/inventory-service/repository"
	"github.com/IBM/sarama"
	"google.golang.org/grpc/metadata"
	"gorm.io/gorm"

	kafkautil "github.com/BhavaniNBL/ecommerce-backend/shared/kafkautil"

	pb "github.com/BhavaniNBL/ecommerce-backend/proto/inventorypb"
	"github.com/BhavaniNBL/ecommerce-backend/proto/productpb"
	"github.com/redis/go-redis/v9"
)

type InventoryService struct {
	pb.UnimplementedInventoryServiceServer
	repo          repository.InventoryRepo
	kafkaProducer sarama.SyncProducer
	productClient productpb.ProductServiceClient
	redis         *redis.Client
	topicName     string
}

func NewInventoryService(kafkaProducer sarama.SyncProducer, redisClient *redis.Client, topicName string, productClient productpb.ProductServiceClient) *InventoryService {
	return &InventoryService{
		kafkaProducer: kafkaProducer,
		redis:         redisClient,
		topicName:     topicName,
		productClient: productClient,
	}
}

func (s *InventoryService) SetRepo(repo repository.InventoryRepo) {
	s.repo = repo
}

func (s *InventoryService) GetInventory(ctx context.Context, req *pb.GetInventoryRequest) (*pb.InventoryResponse, error) {
	// Check Redis cache
	val, err := s.redis.Get(ctx, req.ProductId).Result()
	if err == nil {
		var cached model.Inventory
		if err := json.Unmarshal([]byte(val), &cached); err == nil {
			return &pb.InventoryResponse{
				ProductId:         cached.ProductID,
				Quantity:          cached.Quantity,
				WarehouseLocation: cached.WarehouseLocation,
			}, nil
		}
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Println("⚠️ Failed to get metadata from context")
	}
	outCtx := metadata.NewOutgoingContext(ctx, md)

	productResp, err := s.productClient.CheckProductExists(outCtx, &productpb.ProductID{Id: req.ProductId})

	// Check Product existence
	//productResp, err := s.productClient.CheckProductExists(ctx, &productpb.ProductID{Id: req.ProductId})
	if err != nil {
		log.Println("❌ [InventoryService] Error checking product:", err)
	}
	if productResp != nil {
		log.Println("✅ [InventoryService] Product exists check:", productResp.Exists)
	}
	if err != nil || !productResp.Exists {
		return nil, fmt.Errorf("Product does not exist")
	}
	inv, err := s.repo.GetInventory(req.ProductId)
	if err != nil {
		return nil, err
	}

	// Cache in Redis
	data, _ := json.Marshal(inv)
	s.redis.Set(ctx, req.ProductId, data, 0)

	return &pb.InventoryResponse{
		ProductId:         inv.ProductID,
		Quantity:          inv.Quantity,
		WarehouseLocation: inv.WarehouseLocation,
	}, nil
}

// func (s *InventoryService) UpdateInventory(ctx context.Context, req *pb.UpdateInventoryRequest) (*pb.InventoryResponse, error) {
// 	//productResp, err := s.productClient.CheckProductExists(ctx, &productpb.ProductID{Id: req.ProductId})
// 	md, ok := metadata.FromIncomingContext(ctx)
// 	if !ok {
// 		log.Println("⚠️ Failed to get metadata from context")
// 	}
// 	outCtx := metadata.NewOutgoingContext(ctx, md)

// 	productResp, err := s.productClient.CheckProductExists(outCtx, &productpb.ProductID{Id: req.ProductId})
// 	fmt.Println("📞 Checking product ID:", req.ProductId)
// 	fmt.Println("✅ Product exists response:", productResp, "err:", err)
// 	if err != nil || !productResp.Exists {
// 		return nil, fmt.Errorf("Product does not exist")
// 	}

// 	inv, err := s.repo.UpdateInventory(req.ProductId, req.QuantityChange)
// 	if err != nil {
// 		return nil, err
// 	}

// 	eventType := "InventoryReserved"
// 	if inv.Quantity < 0 {
// 		eventType = "OutOfStock"
// 	}

// 	eventPayload := map[string]interface{}{
// 		"type":       eventType,
// 		"product_id": req.ProductId,
// 		"quantity":   inv.Quantity,
// 	}
// 	data, _ := json.Marshal(eventPayload)

// 	// Using the new sarama publisher
// 	kafkautil.SafePublish("localhost:9092", "inventory-events", req.ProductId, string(data))

// 	s.redis.Del(ctx, req.ProductId) // invalidate cache

// 	return &pb.InventoryResponse{
// 		ProductId:         inv.ProductID,
// 		Quantity:          inv.Quantity,
// 		WarehouseLocation: inv.WarehouseLocation,
// 	}, nil
// }

func (s *InventoryService) UpdateInventory(ctx context.Context, req *pb.UpdateInventoryRequest) (*pb.InventoryResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Println("⚠️ Failed to get metadata from context")
	}
	outCtx := metadata.NewOutgoingContext(ctx, md)

	productResp, err := s.productClient.CheckProductExists(outCtx, &productpb.ProductID{Id: req.ProductId})
	fmt.Println("📞 Checking product ID:", req.ProductId)
	fmt.Println("✅ Product exists response:", productResp, "err:", err)
	if err != nil || !productResp.Exists {
		return nil, fmt.Errorf("Product does not exist")
	}

	// ✅ NEW LOGIC START
	inv, err := s.repo.GetInventory(req.ProductId)
	if err != nil && err == gorm.ErrRecordNotFound {
		newInv := &model.Inventory{
			ProductID:         req.ProductId,
			Quantity:          req.QuantityChange,
			WarehouseLocation: "Default Warehouse",
		}
		if err := s.repo.CreateInventory(newInv); err != nil {
			return nil, err
		}
		inv = newInv
	} else if err == nil {
		inv.Quantity += req.QuantityChange
		if err := s.repo.UpdateExistingInventory(inv); err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}
	// ✅ NEW LOGIC END

	eventType := "InventoryReserved"
	if inv.Quantity < 0 {
		eventType = "OutOfStock"
	}

	eventPayload := map[string]interface{}{
		"type":       eventType,
		"product_id": req.ProductId,
		"quantity":   inv.Quantity,
	}
	data, _ := json.Marshal(eventPayload)

	kafkautil.SafePublish("localhost:9092", "inventory-events", req.ProductId, string(data))
	s.redis.Del(ctx, req.ProductId)

	return &pb.InventoryResponse{
		ProductId:         inv.ProductID,
		Quantity:          inv.Quantity,
		WarehouseLocation: inv.WarehouseLocation,
	}, nil
}
