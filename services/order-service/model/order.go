package model

import "time"

type Order struct {
	ID         string    `json:"id" gorm:"primaryKey"`
	ProductID  string    `json:"product_id"`
	Quantity   int32     `json:"quantity"`
	Status     string    `json:"status"` // pending, confirmed, failed
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
}
