package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"
	"os"

	"github.com/BhavaniNBL/ecommerce-backend/services/notification-service/model"
	"github.com/BhavaniNBL/ecommerce-backend/services/notification-service/repository"
	"github.com/BhavaniNBL/ecommerce-backend/shared/kafkautil"
)

type NotificationService struct {
	Repo        *repository.NotificationRepo
	KafkaBroker string
	Topic       string
}

func NewNotificationService(repo *repository.NotificationRepo, broker, topic string) *NotificationService {
	return &NotificationService{Repo: repo, KafkaBroker: broker, Topic: topic}
}

func (s *NotificationService) HandleOrderConfirmedEvent(message []byte) {
	var event map[string]interface{}
	if err := json.Unmarshal(message, &event); err != nil {
		log.Println("❌ Failed to unmarshal event:", err)
		return
	}
	if event["type"] != "OrderConfirmed" {
		log.Println("⚠️ Ignored event of type:", event["type"])
		return
	}

	notification := &model.Notification{
		OrderID: event["order_id"].(string),
		UserID:  event["user_id"].(string),
		Message: "Your order has been confirmed!",
		Channel: "email",
	}

	if err := s.Repo.Save(notification); err != nil {
		log.Println("❌ Failed to store notification:", err)
		return
	}

	s.Send(notification)
}

func (s *NotificationService) Send(n *model.Notification) {
	switch n.Channel {
	case "email":
		s.sendSMTP(n)
	case "sms":
		log.Printf("📲 Simulated SMS sent to user %s for order %s\n", n.UserID, n.OrderID)
	default:
		log.Printf("📨 Notification (channel: %s) sent to user %s for order %s\n", n.Channel, n.UserID, n.OrderID)
	}
}

func (s *NotificationService) sendSMTP(n *model.Notification) {
	from := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	to := []string{n.UserID + "@gmail.com"}
	subject := "Subject: Notification\n"
	body := fmt.Sprintf("Hi!\n\n%s\nOrder ID: %s\n", n.Message, n.OrderID)

	msg := []byte(subject + "\n" + body)
	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, msg)
	if err != nil {
		log.Println("❌ Failed to send email:", err)
		alert := map[string]interface{}{
			"type":     "NotificationFailed",
			"order_id": n.OrderID,
			"user_id":  n.UserID,
			"channel":  "email",
			"error":    err.Error(),
		}
		data, _ := json.Marshal(alert)
		kafkautil.SafePublish(s.KafkaBroker, s.Topic, n.OrderID, string(data))
	} else {
		log.Printf("📧 Email sent to %s for order %s\n", to[0], n.OrderID)
	}
}
