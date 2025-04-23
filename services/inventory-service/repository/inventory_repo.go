package repository

import (
	"fmt"
	"net/http"

	"github.com/BhavaniNBL/ecommerce-backend/services/inventory-service/model"
	"gorm.io/gorm"
)

type InventoryRepo interface {
	GetInventory(productID string) (*model.Inventory, error)
	UpdateInventory(productID string, change int32) (*model.Inventory, error)
	CreateInventory(inv *model.Inventory) error
	UpdateExistingInventory(inv *model.Inventory) error
}

type inventoryRepo struct {
	db *gorm.DB
}

func NewInventoryRepo(db *gorm.DB) InventoryRepo {
	return &inventoryRepo{db: db}
}

func (r *inventoryRepo) CreateInventory(inv *model.Inventory) error {
	return r.db.Create(inv).Error
}

func (r *inventoryRepo) UpdateExistingInventory(inv *model.Inventory) error {
	return r.db.Save(inv).Error
}

func (r *inventoryRepo) GetInventory(productID string) (*model.Inventory, error) {
	var inv model.Inventory
	if err := r.db.First(&inv, "product_id = ?", productID).Error; err != nil {
		return nil, err
	}
	return &inv, nil
}

func (r *inventoryRepo) UpdateInventory(productID string, change int32) (*model.Inventory, error) {
	if !productExists(productID) {
		return nil, fmt.Errorf("Product not found")
	}
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

func productExists(productID string) bool {
	url := fmt.Sprintf("http://localhost:8082/products/%s", productID)
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != 200 {
		return false
	}
	return true
}
