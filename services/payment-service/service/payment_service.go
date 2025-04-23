package service

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"github.com/BhavaniNBL/ecommerce-backend/services/payment-service/model"
	"github.com/BhavaniNBL/ecommerce-backend/services/payment-service/repository"
	"github.com/BhavaniNBL/ecommerce-backend/shared/kafkautil"
	"github.com/google/uuid"
)

type PaymentService struct {
	KafkaBroker string
	InputTopic  string
	OutputTopic string
	Repo        *repository.PaymentRepo
}

func NewPaymentService(broker, inputTopic, outputTopic string, repo *repository.PaymentRepo) *PaymentService {
	return &PaymentService{
		KafkaBroker: broker,
		InputTopic:  inputTopic,
		OutputTopic: outputTopic,
		Repo:        repo,
	}
}

func (s *PaymentService) ProcessOrderEvent(message []byte) {
	var event map[string]interface{}
	if err := json.Unmarshal(message, &event); err != nil {
		log.Println("❌ Failed to parse order event:", err)
		return
	}

	if event["type"] != "OrderCreated" {
		log.Println("⚠️ Ignored non-OrderCreated event")
		return
	}

	orderID := event["order_id"].(string)
	productID := event["product_id"].(string)
	quantity := int32(event["quantity"].(float64))

	status := model.StatusSuccess
	if rand.Intn(10) < 2 {
		status = model.StatusFailed
	}

	record := &model.Payment{
		ID:        uuid.New().String(),
		OrderID:   orderID,
		ProductID: productID,
		Quantity:  quantity,
		Status:    string(status),
		Timestamp: time.Now(),
	}
	_ = s.Repo.Save(record)

	paymentEvent := model.PaymentEvent{
		Type:      "PaymentProcessed",
		OrderID:   orderID,
		ProductID: productID,
		Quantity:  quantity,
		Status:    string(status),
	}

	data, _ := json.Marshal(paymentEvent)
	kafkautil.SafePublish(s.KafkaBroker, s.OutputTopic, orderID, string(data))
	log.Printf("💸 Payment %s for order %s\n", status, orderID)
}
