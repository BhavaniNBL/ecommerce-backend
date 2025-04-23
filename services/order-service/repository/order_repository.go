package repository

import (
	"github.com/BhavaniNBL/ecommerce-backend/services/order-service/model"
	"gorm.io/gorm"
)

type OrderRepository struct {
	DB *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{DB: db}
}

func (r *OrderRepository) CreateOrder(order *model.Order) error {
	return r.DB.Create(order).Error
}

func (r *OrderRepository) UpdateStatus(id string, status string) error {
	return r.DB.Model(&model.Order{}).Where("id = ?", id).Update("status", status).Error
}

func (r *OrderRepository) GetOrderByID(id string) (*model.Order, error) {
	var order model.Order
	err := r.DB.First(&order, "id = ?", id).Error
	return &order, err
}

func (r *OrderRepository) ListOrders() ([]model.Order, error) {
	var orders []model.Order
	err := r.DB.Find(&orders).Error
	return orders, err
}
