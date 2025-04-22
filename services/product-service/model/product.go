// package model

// import (
// 	"time"

// 	"github.com/google/uuid"
// 	"gorm.io/gorm"
// )

// type Product struct {
// 	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
// 	Name        string    `gorm:"not null" json:"name"`
// 	Description string    `json:"description"`
// 	Price       float64   `gorm:"not null" json:"price"`
// 	Category    string    `json:"category"`
// 	Stock       int       `json:"stock"`
// 	CreatedAt   time.Time `json:"created_at"`
// 	UpdatedAt   time.Time `json:"updated_at"`
// }

// func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {
// 	p.ID = uuid.New()
// 	return
// }

package model

import (
	"time"

	"github.com/BhavaniNBL/ecommerce-backend/proto/productpb" // Adjust the import path for productpb
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`
	Price       float32   `gorm:"not null" json:"price"`
	Category    string    `json:"category"`
	Quantity    int32     `gorm:"not null" json:"quantity"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID = uuid.New() // Ensure a new UUID is assigned if it's not already set
	}
	return
}

// Convert Product to gRPC Product (including Timestamp conversion)
func (p *Product) ToGRPC() *productpb.Product {
	return &productpb.Product{
		Id:          p.ID.String(),
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Category:    p.Category,
		Quantity:    p.Quantity,
		CreatedAt:   timestamppb.New(p.CreatedAt),
		UpdatedAt:   timestamppb.New(p.UpdatedAt),
	}
}

// // Helper function to convert time.Time to google.protobuf.Timestamp
// func toTimestamp(t time.Time) *productpb.Timestamp {
// 	if t.IsZero() {
// 		return nil
// 	}
// 	return &productpb.Timestamp{
// 		Seconds: t.Unix(),
// 		Nanos:   int32(t.Nanosecond()),
// 	}
// }
