package servicebus

import (
	"context"
	"log"

	servicebus "github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

type Receiver struct {
	receiver *servicebus.Receiver
}

func NewReceiver(connectionString, queueName string) (*Receiver, error) {
	client, err := servicebus.NewClientFromConnectionString(connectionString, &servicebus.ClientOptions{})
	if err != nil {
		return nil, err
	}

	receiver, err := client.NewReceiverForQueue(queueName, &servicebus.ReceiverOptions{})
	if err != nil {
		return nil, err
	}

	return &Receiver{
		receiver: receiver,
	}, nil
}

func (r *Receiver) ReceiveMessages(ctx context.Context) error {

	for {
		messages, err := r.receiver.ReceiveMessages(ctx, 1, &servicebus.ReceiveMessagesOptions{})
		if err != nil {
			return err
		}

		for _, message := range messages {
			log.Printf("Received message: %s", string(message.Body))
			if err := r.receiver.CompleteMessage(ctx, message, &servicebus.CompleteMessageOptions{}); err != nil {
				if err = r.receiver.AbandonMessage(ctx, message, &servicebus.AbandonMessageOptions{}); err != nil {
					return err
				}
				return err
			}
		}
	}
}

func (s *Receiver) Close() error {
	return s.receiver.Close(context.Background())
}
