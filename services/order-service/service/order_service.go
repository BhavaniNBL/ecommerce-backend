package service

import (
	"context"
	"encoding/json"
	"log"
	"time"

	inventorypb "github.com/BhavaniNBL/ecommerce-backend/proto/inventorypb"
	"github.com/BhavaniNBL/ecommerce-backend/services/order-service/model"
	"github.com/BhavaniNBL/ecommerce-backend/services/order-service/repository"
	"github.com/BhavaniNBL/ecommerce-backend/shared/kafkautil"
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

type OrderService struct {
	Repo            *repository.OrderRepository
	InventoryClient inventorypb.InventoryServiceClient
	KafkaBroker     string
	KafkaTopic      string
}

func NewOrderService(repo *repository.OrderRepository, inv inventorypb.InventoryServiceClient, broker, topic string) *OrderService {
	return &OrderService{
		Repo:            repo,
		InventoryClient: inv,
		KafkaBroker:     broker,
		KafkaTopic:      topic,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, productID string, qty int32) error {
	id := uuid.New().String()

	//_, err := s.InventoryClient.UpdateInventory(ctx, &inventorypb.UpdateInventoryRequest{
	//	ProductId:      productID,
	//	QuantityChange: -qty,
	//})
	// 🔐 Inject auth metadata from incoming context (HTTP) into outgoing gRPC context
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		log.Println("📦 Forwarding metadata:", md)
		ctx = metadata.NewOutgoingContext(ctx, md)
	} else {
		log.Println("❌ No metadata found in incoming context")
	}

	_, err := s.InventoryClient.UpdateInventory(ctx, &inventorypb.UpdateInventoryRequest{
		ProductId:      productID,
		QuantityChange: -qty,
	})
	if err != nil {
		return err
	}

	order := &model.Order{
		ID:         id,
		ProductID:  productID,
		Quantity:   qty,
		Status:     "pending",
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}
	if err := s.Repo.CreateOrder(order); err != nil {
		return err
	}

	event := map[string]interface{}{
		"type":       "OrderCreated",
		"order_id":   id,
		"product_id": productID,
		"quantity":   qty,
	}
	payload, _ := json.Marshal(event)
	kafkautil.SafePublish(s.KafkaBroker, s.KafkaTopic, id, string(payload))

	return nil
}

func (s *OrderService) GetOrderByID(id string) (*model.Order, error) {
	return s.Repo.GetOrderByID(id)
}

func (s *OrderService) ListOrders() ([]model.Order, error) {
	return s.Repo.ListOrders()
}

func (s *OrderService) UpdateOrderStatus(id, status string) error {
	return s.Repo.UpdateStatus(id, status)
}
