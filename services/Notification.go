package services

import (
	"dirwatcher/models"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

type NotificationService struct {
	notifications []models.Notification
	idCounter     int
}

func NewNotificationService() *NotificationService {
	return &NotificationService{
		notifications: []models.Notification{},
		idCounter:     1,
	}
}

func (ns *NotificationService) SendNotification(userID int, message string) models.Notification {
	notification := models.Notification{
		ID:      ns.idCounter,
		UserID:  userID,
		Message: message,
		Status:  "sent",
	}
	ns.notifications = append(ns.notifications, notification)
	ns.idCounter++
	return notification
}

func (ns *NotificationService) GetNotifications(userID int) []models.Notification {
	var userNotifications []models.Notification
	for _, notification := range ns.notifications {
		if notification.UserID == userID {
			userNotifications = append(userNotifications, notification)
		}
	}
	return userNotifications
}

func sendFCMNotification(message string) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"), // Change to your desired region
	})
	if err != nil {
		log.Fatalf("failed to create session: %v", err)
	}

	// Create SNS service client
	svc := sns.New(sess)

	// Define your SNS topic ARN
	topicARN := "arn:aws:sns:us-east-1:YOUR_ACCOUNT_ID:YOUR_TOPIC_NAME"

	// Define the message to be sent
	input := &sns.PublishInput{
		Message:  aws.String(message),
		TopicArn: aws.String(topicARN),
	}

	result, err := svc.Publish(input)
	if err != nil {
		return err
	}

	fmt.Printf("Message sent to the topic %s. Message ID: %s\n", topicARN, *result.MessageId)
	return nil

}
