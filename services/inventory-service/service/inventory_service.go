package service

import (
	"context"
	"encoding/json"

	"github.com/BhavaniNBL/ecommerce-backend/services/inventory-service/model"
	"github.com/BhavaniNBL/ecommerce-backend/services/inventory-service/repository"
	"github.com/BhavaniNBL/ecommerce-backend/shared/kafka"

	pb "github.com/BhavaniNBL/ecommerce-backend/proto/inventorypb"
	"github.com/redis/go-redis/v9"
)

type InventoryService struct {
	pb.UnimplementedInventoryServiceServer
	repo        repository.InventoryRepo
	kafkaWriter *kafka.Writer
	redis       *redis.Client
	topicName   string
}

func NewInventoryService(kafkaProducer *kafka.Writer, redisClient *redis.Client, topicName string) *InventoryService {
	return &InventoryService{
		kafkaWriter: kafkaProducer,
		redis:       redisClient,
		topicName:   topicName,
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

func (s *InventoryService) UpdateInventory(ctx context.Context, req *pb.UpdateInventoryRequest) (*pb.InventoryResponse, error) {
	inv, err := s.repo.UpdateInventory(req.ProductId, req.QuantityChange)
	if err != nil {
		return nil, err
	}

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
	kafka.PublishMessage(s.kafkaWriter, "inventory-events", req.ProductId, string(data))

	s.redis.Del(ctx, req.ProductId) // invalidate cache

	return &pb.InventoryResponse{
		ProductId:         inv.ProductID,
		Quantity:          inv.Quantity,
		WarehouseLocation: inv.WarehouseLocation,
	}, nil
}
