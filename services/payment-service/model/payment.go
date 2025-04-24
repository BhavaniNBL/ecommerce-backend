package model

import "time"

type PaymentStatus string

const (
	StatusSuccess PaymentStatus = "success"
	StatusFailed  PaymentStatus = "failed"
)

// PaymentEvent represents the structure of the payment event
/* This struct is used to publish events to Kafka, like "PaymentProcessed".

These events are consumed by other services (e.g., Order or Notification service). */
type PaymentEvent struct {
	Type      string `json:"type"`
	OrderID   string `json:"order_id"`
	ProductID string `json:"product_id"`
	Quantity  int32  `json:"quantity"`
	Status    string `json:"status"`
}

/* This is the main database model used with GORM:
This struct is mapped to the payments table in Postgres automatically via db.AutoMigrate() */

type Payment struct {
	ID        string `gorm:"primaryKey"`
	OrderID   string
	ProductID string
	Quantity  int32
	Status    string
	Amount    float64
	Timestamp time.Time
}
