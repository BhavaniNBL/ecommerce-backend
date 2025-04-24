package model

type Notification struct {
	OrderID string `json:"order_id"`
	UserID  string `json:"user_id"`
	Message string `json:"message"`
	Channel string `json:"channel"` // email, sms, push
}
