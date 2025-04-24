package handler

import (
	"net/http"

	"github.com/BhavaniNBL/ecommerce-backend/services/notification-service/model"
	"github.com/BhavaniNBL/ecommerce-backend/services/notification-service/repository"
	"github.com/BhavaniNBL/ecommerce-backend/services/notification-service/service"
	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	Repo  *repository.NotificationRepo
	Logic *service.NotificationService
}

func NewNotificationHandler(repo *repository.NotificationRepo, logic *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{Repo: repo, Logic: logic}
}

func (h *NotificationHandler) GetByOrderID(c *gin.Context) {
	orderID := c.Param("orderID")
	notif, err := h.Repo.Get(orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "notification not found"})
		return
	}
	c.JSON(http.StatusOK, notif)
}

func (h *NotificationHandler) CreateNotification(c *gin.Context) {
	var req model.Notification
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.Repo.Save(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save notification"})
		return
	}
	h.Logic.Send(&req)
	c.JSON(http.StatusOK, gin.H{"message": "notification created and sent"})
}
