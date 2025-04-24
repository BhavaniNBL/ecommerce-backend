package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/BhavaniNBL/ecommerce-backend/services/notification-service/model"
	"github.com/redis/go-redis/v9"
)

type NotificationRepo struct {
	RedisClient *redis.Client
	Ctx         context.Context
}

func NewNotificationRepo(rdb *redis.Client) *NotificationRepo {
	return &NotificationRepo{
		RedisClient: rdb,
		Ctx:         context.Background(),
	}
}

func (r *NotificationRepo) Save(notification *model.Notification) error {
	key := fmt.Sprintf("notification:%s", notification.OrderID)
	data, err := json.Marshal(notification)
	if err != nil {
		return err
	}
	return r.RedisClient.Set(r.Ctx, key, data, 0).Err()
}

func (r *NotificationRepo) Get(orderID string) (*model.Notification, error) {
	key := fmt.Sprintf("notification:%s", orderID)
	val, err := r.RedisClient.Get(r.Ctx, key).Result()
	if err != nil {
		return nil, err
	}
	var notif model.Notification
	err = json.Unmarshal([]byte(val), &notif)
	return &notif, err
}
