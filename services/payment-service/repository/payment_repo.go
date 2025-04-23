package repository

import (
	"github.com/BhavaniNBL/ecommerce-backend/services/payment-service/model"
	"gorm.io/gorm"
)

type PaymentRepo struct {
	DB *gorm.DB
}

func NewPaymentRepo(db *gorm.DB) *PaymentRepo {
	return &PaymentRepo{DB: db}
}

func (r *PaymentRepo) Save(payment *model.Payment) error {
	return r.DB.Create(payment).Error
}

func (r *PaymentRepo) FindByOrderID(orderID string) ([]model.Payment, error) {
	var records []model.Payment
	err := r.DB.Where("order_id = ?", orderID).Find(&records).Error
	return records, err
}
