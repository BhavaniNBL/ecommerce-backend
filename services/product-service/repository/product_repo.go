package repository

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	config "github.com/BhavaniNBL/ecommerce-backend/config/redis"
	"github.com/BhavaniNBL/ecommerce-backend/services/product-service/model"
	"github.com/google/uuid"

	"gorm.io/gorm"
)

type ProductRepository struct {
	DB *gorm.DB
}

func NewProductRepo(db *gorm.DB) *ProductRepository {
	return &ProductRepository{DB: db}
}

func (r *ProductRepository) Create(product *model.Product) error {
	return r.DB.Create(product).Error
}

// func (r *ProductRepository) GetByID(id string) (*model.Product, error) {
// 	cacheKey := fmt.Sprintf("product:%s", id)
// 	val, err := config.RedisClient.Get(config.Ctx, cacheKey).Result()
// 	if err == nil {
// 		var prod model.Product
// 		if err := json.Unmarshal([]byte(val), &prod); err == nil {
// 			return &prod, nil
// 		}
// 	}

// 	// Parse UUID properly
// 	parsedID, err := uuid.Parse(id)
// 	if err != nil {
// 		log.Println("UUID parsing failed:", err)
// 		return nil, err
// 	}

// 	var product model.Product
// 	if err := r.DB.First(&product, "id = ?", parsedID).Error; err != nil {
// 		log.Println("DB fetch failed:", err)
// 		return nil, err
// 	}

// 	// Cache it
// 	jsonData, _ := json.Marshal(product)
// 	config.RedisClient.Set(config.Ctx, cacheKey, jsonData, 10*time.Minute)

// 	return &product, nil
// }

func (r *ProductRepository) GetByID(id string) (*model.Product, error) {
	log.Println("GetByID called with ID:", id)

	// 1. Validate UUID
	parsedID, err := uuid.Parse(id)
	if err != nil {
		log.Println("UUID parsing failed:", err)
		return nil, fmt.Errorf("invalid UUID format: %w", err)
	}

	// 2. Try Redis first
	cacheKey := fmt.Sprintf("product:%s", parsedID.String())
	val, err := config.RedisClient.Get(config.Ctx, cacheKey).Result()
	if err == nil {
		var prod model.Product
		if err := json.Unmarshal([]byte(val), &prod); err == nil {
			log.Println("Cache hit, returning product from Redis")
			return &prod, nil
		}
		log.Println("Cache unmarshal failed:", err)
	}

	// 3. Query DB
	var product model.Product
	if err := r.DB.First(&product, "id = ?", parsedID).Error; err != nil {
		log.Println("DB fetch failed:", err)
		return nil, fmt.Errorf("database error: %w", err)
	}

	// 4. Cache it for future
	jsonData, _ := json.Marshal(product)
	err = config.RedisClient.Set(config.Ctx, cacheKey, jsonData, 10*time.Minute).Err()
	if err != nil {
		log.Println("Redis set failed:", err)
	}

	log.Println("Returning product from DB")
	return &product, nil
}

func (r *ProductRepository) Update(product *model.Product) error {
	return r.DB.Save(product).Error
}

func (r *ProductRepository) Delete(id string) error {
	config.RedisClient.Del(config.Ctx, fmt.Sprintf("product:%s", id))
	return r.DB.Delete(&model.Product{}, "id = ?", id).Error
}

func (r *ProductRepository) List(filters map[string]string) ([]model.Product, error) {
	query := r.DB.Model(&model.Product{})
	if name, ok := filters["name"]; ok {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}
	if cat, ok := filters["category"]; ok {
		query = query.Where("category = ?", cat)
	}
	var products []model.Product
	err := query.Find(&products).Error
	return products, err
}
