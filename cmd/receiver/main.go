package main

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"time"

	servicebus "github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

func main() {
	// Replace with your Service Bus connection string and queue name
	connectionString := "<service bus connection string>" //os.Getenv("SERVICE_BUS_CONNECTION_STRING")
	queueName := "<queue name>"                           //os.Getenv("QUEUE_NAME")

	if connectionString == "" || queueName == "" {
		log.Fatal("Please set the SERVICE_BUS_CONNECTION_STRING and QUEUE_NAME environment variables.")
	}

	// Create a new Service Bus client
	sbClient, err := servicebus.NewClientFromConnectionString(connectionString, &servicebus.ClientOptions{})
	if err != nil {
		log.Fatalf("Failed to create Service Bus client: %v", err)
	}

	// Create a new receiver for the queue
	receiver, err := sbClient.NewReceiverForQueue(queueName, &servicebus.ReceiverOptions{})
	if err != nil {
		log.Fatalf("Failed to create receiver: %v", err)
	}

	// Start receiving messages, this will block until the context is cancelled
	// You can use a context with a timeout to limit the time spent waiting for messages
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	defer receiver.Close(ctx)
	defer sbClient.Close(ctx)

	log.Println("Waiting to receive message... queue reception will stopped after 60 seconds")
	messageOptions := &servicebus.ReceiveMessagesOptions{}
	rand.New(rand.NewSource(time.Now().UnixNano())) // Seed the random number generator

	for {
		messages, err := receiver.ReceiveMessages(ctx, 1, messageOptions)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				log.Println("Context deadline exceeded. Stopping message reception.")
				break
			}

			log.Printf("Failed to receive message: %v", err)
			continue
		}

		log.Println("New message incoming!")
		for _, message := range messages {
			// Process the message
			log.Printf("Received message: %s\n", string(message.Body))

			stopRenew := make(chan bool) // Channel to signal the Goroutine to stop
			// Renew the lock periodically
			ticker := time.NewTicker(10 * time.Second) // Renew every 10 seconds. Message lock is 60 seconds by default
			defer ticker.Stop()

			go func(msg *servicebus.ReceivedMessage, ch chan bool) {
				for {
					select {
					case <-ch:
						log.Println("Kill renew message goroutine")
						return
					case <-ticker.C:
						if err := receiver.RenewMessageLock(ctx, msg, &servicebus.RenewMessageLockOptions{}); err != nil {
							log.Printf("Failed to renew message lock: %s", err)
						} else {
							log.Println("Message lock renewed.")
						}
					}
				}
			}(message, stopRenew)

			// Process message here
			log.Printf("Processing message for 15 seconds")
			time.Sleep(15 * time.Second) //sleep for 15 seconds to simulate processing

			// Complete the message so it is not received again
			if err := receiver.CompleteMessage(context.TODO(), message, &servicebus.CompleteMessageOptions{}); err != nil {
				var sbErr *servicebus.Error

				if errors.As(err, &sbErr) && sbErr.Code == servicebus.CodeLockLost {
					// The message lock has expired. This isn't fatal for the client, but it does mean
					// that this message can be received by another Receiver (or potentially this one).
					log.Printf("Message lock expired")

					// Alternatively, you can extend the message lock by calling Receiver.RenewMessageLock(msg) before the
					// message lock has expired.
					continue
				}

				// Other errors when can't complete message
				if err := receiver.AbandonMessage(context.TODO(), message, &servicebus.AbandonMessageOptions{}); err != nil {
					log.Printf("Failed to unlock message for retry: %v", err)
				} else {
					log.Printf("Message abandoned and put to queue")
				}

				log.Printf("Failed to complete message: %v", err)
				stopRenew <- true
			} else {
				log.Printf("Message completed: %s\n", string(message.Body))
				stopRenew <- true
			}

			sleepDuration := time.Duration(rand.Intn(5)+1) * time.Second // Random sleep between 1 and 5 seconds
			log.Printf("Sleeping for %v before processing the next message...", sleepDuration)
			time.Sleep(sleepDuration)
		}
	}
}
