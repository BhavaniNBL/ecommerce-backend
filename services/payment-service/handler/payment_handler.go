package handler

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"github.com/BhavaniNBL/ecommerce-backend/services/payment-service/model"
	"github.com/BhavaniNBL/ecommerce-backend/services/payment-service/repository"
	"github.com/BhavaniNBL/ecommerce-backend/shared/kafkautil"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PaymentHandler struct {
	Repo *repository.PaymentRepo
}

func NewPaymentHandler(repo *repository.PaymentRepo) *PaymentHandler {
	return &PaymentHandler{Repo: repo}
}

type PaymentRequest struct {
	OrderID    string  `json:"order_id" binding:"required"`
	ProductID  string  `json:"product_id" binding:"required"`
	Quantity   int32   `json:"quantity" binding:"required"`
	Amount     float64 `json:"amount" binding:"required"`
	CardNumber string  `json:"card_number" binding:"required"`
	Expiry     string  `json:"expiry" binding:"required"`
	CVV        string  `json:"cvv" binding:"required"`
}

func (h *PaymentHandler) ProcessPayment(c *gin.Context) {
	var req PaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(req.CardNumber) != 16 || req.CVV == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid card details"})
		return
	}

	status := model.StatusSuccess
	if rand.Intn(10) < 2 {
		status = model.StatusFailed
	}

	record := &model.Payment{
		ID:        uuid.New().String(),
		OrderID:   req.OrderID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
		Amount:    req.Amount,
		Status:    string(status),
		Timestamp: time.Now(),
	}
	h.Repo.Save(record)

	paymentEvent := model.PaymentEvent{
		Type:      "PaymentProcessed",
		OrderID:   req.OrderID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
		Status:    string(status),
	}

	data, _ := json.Marshal(paymentEvent)
	kafkautil.SafePublish(h.Repo.KafkaBroker, "payment-events", req.OrderID, string(data))

	c.JSON(http.StatusOK, gin.H{"message": "Payment processed", "status": status})
}

func (h *PaymentHandler) GetByOrderID(c *gin.Context) {
	orderID := c.Param("orderID")
	records, err := h.Repo.FindByOrderID(orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch records"})
		return
	}

	if len(records) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no payments found for this order"})
		return
	}
	c.JSON(http.StatusOK, records)
}
