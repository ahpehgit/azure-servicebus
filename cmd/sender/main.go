package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	servicebus "github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

func main() {
	// Replace with your Azure Service Bus connection string and queue name
	connectionString := "<service bus connection string>" //os.Getenv("SERVICE_BUS_CONNECTION_STRING")
	queueName := "<queue name>"                           //os.Getenv("QUEUE_NAME")

	// Create a new Service Bus client
	sbClient, err := servicebus.NewClientFromConnectionString(connectionString, nil)
	if err != nil {
		log.Fatalf("failed to create Service Bus client: %v", err)
	}

	// Create a new sender for the queue
	sender, err := sbClient.NewSender(queueName, nil)
	if err != nil {
		log.Fatalf("failed to create sender: %v", err)
	}

	messageOptions := &servicebus.SendMessageOptions{}

	// Send the message
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	defer sender.Close(ctx)
	defer sbClient.Close(ctx)

	var wg sync.WaitGroup
	for i := range 10 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			// Prepare a message to send
			message := &servicebus.Message{
				Body: []byte(fmt.Sprintf("Hello, Azure Service Bus! This is message %d!", i)),
			}

			err = sender.SendMessage(ctx, message, messageOptions)
			if err != nil {
				log.Fatalf("failed to send message: %v", err)
			}

			log.Println("Message sent successfully!")
		}(i)
	}
	wg.Wait() // Wait for all goroutines to complete
	log.Println("All messages sent!")
}
