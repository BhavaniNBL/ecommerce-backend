package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/BhavaniNBL/ecommerce-backend/services/order-service/service"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
)

type OrderHandler struct {
	Service *service.OrderService
}

func NewOrderHandler(service *service.OrderService) *OrderHandler {
	return &OrderHandler{Service: service}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req struct {
		ProductID string `json:"product_id"`
		Quantity  int32  `json:"quantity"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token := c.GetHeader("Authorization")
	log.Println("🔐 Token received:", token)

	ctx := metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{
		"authorization": c.GetHeader("Authorization"),
	}))
	//err := h.Service.CreateOrder(context.Background(), req.ProductID, req.Quantity)
	err := h.Service.CreateOrder(ctx, req.ProductID, req.Quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order created successfully"})
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	id := c.Param("id")
	order, err := h.Service.GetOrderByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}
	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) ListOrders(c *gin.Context) {
	orders, err := h.Service.ListOrders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.Service.UpdateOrderStatus(id, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Order status updated"})
}
