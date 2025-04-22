package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/BhavaniNBL/ecommerce-backend/services/inventory-service/model"
	"github.com/redis/go-redis/v9"
)

type Cache interface {
	GetInventory(productID string) (*model.Inventory, error)
	SetInventory(productID string, data *model.Inventory) error
}

type redisCache struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisCache(addr string) Cache {
	rdb := redis.NewClient(&redis.Options{Addr: addr})
	return &redisCache{client: rdb, ctx: context.Background()}
}

func (c *redisCache) GetInventory(productID string) (*model.Inventory, error) {
	val, err := c.client.Get(c.ctx, productID).Result()
	if err != nil {
		return nil, err
	}

	var inv model.Inventory
	if err := json.Unmarshal([]byte(val), &inv); err != nil {
		return nil, err
	}
	return &inv, nil
}

func (c *redisCache) SetInventory(productID string, data *model.Inventory) error {
	val, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return c.client.Set(c.ctx, productID, val, time.Minute*10).Err()
}
