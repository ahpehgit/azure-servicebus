# Azure Service Bus Application

This project demonstrates how to use Azure Service Bus with Go to send and receive messages. It consists of two main components: a sender application that sends messages to a queue and a receiver application that listens for and processes those messages.

## Project Structure

```
azure-service-bus-app
├── cmd
│   ├── sender
│   │   └── main.go       # Entry point for the sender application
│   └── receiver
│       └── main.go       # Entry point for the receiver application
├── internal
│   ├── servicebus
│   │   ├── sender.go     # Contains the Sender struct and methods
│   │   └── receiver.go    # Contains the Receiver struct and methods
├── go.mod                 # Module definition for the Go project
├── go.sum                 # Checksums for module dependencies
└── README.md              # Documentation for the project
```

## Prerequisites

- Go 1.23.6 or later
- An Azure account with an active subscription
- Azure Service Bus namespace and queue created

## Setup Instructions

1. Clone the repository:

   ```
   git clone https://github.com/microsoft/vscode-remote-try-dab.git
   cd azure-service-bus-app
   ```

2. Install the necessary dependencies:

   ```
   go mod tidy
   ```

3. Set up your Azure Service Bus connection string and queue name in the environment variables:

   ```
   export SERVICE_BUS_CONNECTION_STRING="your_connection_string"
   export QUEUE_NAME="your_queue_name"
   ```

## Usage

### Sending Messages

To send messages to the Azure Service Bus queue, run the sender application:

```
go run cmd/sender/main.go
```

### Receiving Messages

To receive messages from the Azure Service Bus queue, run the receiver application:

```
go run cmd/receiver/main.go
```

## License

This project is licensed under the MIT License. See the LICENSE file for more details.