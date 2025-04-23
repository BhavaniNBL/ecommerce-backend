package model

import "time"

type PaymentStatus string

const (
	StatusSuccess PaymentStatus = "success"
	StatusFailed  PaymentStatus = "failed"
)

// PaymentEvent represents the structure of the payment event
type PaymentEvent struct {
	Type      string `json:"type"`
	OrderID   string `json:"order_id"`
	ProductID string `json:"product_id"`
	Quantity  int32  `json:"quantity"`
	Status    string `json:"status"`
}

type Payment struct {
	ID        string `gorm:"primaryKey"`
	OrderID   string
	ProductID string
	Quantity  int32
	Status    string
	Timestamp time.Time
}
