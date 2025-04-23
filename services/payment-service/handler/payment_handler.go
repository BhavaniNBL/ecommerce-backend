package handler

import (
	"net/http"

	"github.com/BhavaniNBL/ecommerce-backend/services/payment-service/repository"
	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	Repo *repository.PaymentRepo
}

func NewPaymentHandler(repo *repository.PaymentRepo) *PaymentHandler {
	return &PaymentHandler{Repo: repo}
}

func (h *PaymentHandler) GetByOrderID(c *gin.Context) {
	orderID := c.Param("orderID")
	records, err := h.Repo.FindByOrderID(orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch records"})
		return
	}
	c.JSON(http.StatusOK, records)
}
