package model

import (
	"time"
)

type User struct {
	ID        string    `gorm:"primaryKey;type:uuid" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Email     string    `gorm:"unique;not null" json:"email"`
	Password  string    `gorm:"not null" json:"-"`
	UserType  string    `gorm:"not null" json:"user_type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
// 	u.ID = uuid.New() // Generate a new UUID before saving
// 	return
// }

type SignUpRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	UserType string `json:"user_type" binding:"required"` // Added UserType here
}

// Used in Login handler
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
