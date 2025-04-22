package repository

import (
	"github.com/BhavaniNBL/ecommerce-backend/services/inventory-service/model"
	"gorm.io/gorm"
)

type InventoryRepo interface {
	GetInventory(productID string) (*model.Inventory, error)
	UpdateInventory(productID string, change int32) (*model.Inventory, error)
}

type inventoryRepo struct {
	db *gorm.DB
}

func NewInventoryRepo(db *gorm.DB) InventoryRepo {
	return &inventoryRepo{db: db}
}

func (r *inventoryRepo) GetInventory(productID string) (*model.Inventory, error) {
	var inv model.Inventory
	if err := r.db.First(&inv, "product_id = ?", productID).Error; err != nil {
		return nil, err
	}
	return &inv, nil
}

func (r *inventoryRepo) UpdateInventory(productID string, change int32) (*model.Inventory, error) {
	var inv model.Inventory
	if err := r.db.First(&inv, "product_id = ?", productID).Error; err != nil {
		return nil, err
	}

	inv.Quantity += change
	if err := r.db.Save(&inv).Error; err != nil {
		return nil, err
	}
	return &inv, nil
}
