package repository

import (
	"errors"
	"log"

	"github.com/BhavaniNBL/ecommerce-backend/config/db"
	"github.com/BhavaniNBL/ecommerce-backend/services/user-service/model"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

// SaveUser inserts a new user into the DB
func (r *UserRepository) CreateUser(user *model.User) error {
	result := db.DB.Create(user)
	return result.Error
}

// GetUserByEmail retrieves a user by their email
func (r *UserRepository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	log.Printf("Looking for user by email: %s", email)

	result := db.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("user not found")
	}

	return &user, result.Error
}

// GetUserByID retrieves a user by ID
func (r *UserRepository) GetUserByID(id string) (*model.User, error) {
	var user model.User
	if err := db.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// ListUsers retrieves all users
func (r *UserRepository) ListUsers() ([]model.User, error) {
	var users []model.User
	result := db.DB.Find(&users)
	return users, result.Error
}
