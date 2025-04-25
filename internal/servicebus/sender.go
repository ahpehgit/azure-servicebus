package servicebus

import (
	"context"
	"log"

	servicebus "github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

type Sender struct {
	sender *servicebus.Sender
}

func NewSender(connectionString, queueName string) (*Sender, error) {
	client, err := servicebus.NewClientFromConnectionString(connectionString, &servicebus.ClientOptions{})
	if err != nil {
		return nil, err
	}

	queueSender, err := client.NewSender(queueName, &servicebus.NewSenderOptions{})
	if err != nil {
		return nil, err
	}

	return &Sender{sender: queueSender}, nil
}

func (s *Sender) SendMessage(ctx context.Context, message string) error {
	msg := &servicebus.Message{
		Body: []byte(message),
	}

	err := s.sender.SendMessage(ctx, msg, &servicebus.SendMessageOptions{})
	if err != nil {
		log.Printf("failed to send message: %v", err)
		return err
	}

	log.Println("message sent successfully")
	return nil
}

func (s *Sender) Close() error {
	return s.sender.Close(context.Background())
}
