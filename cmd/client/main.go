package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	connStr := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connStr)
	if err != nil {
		log.Fatalf("Error connecting to RabbitMQ server: %v", err)
	}
	defer conn.Close()

	log.Println("Starting Peril client...")

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	_, queue, err := pubsub.DeclareAndBind(conn, routing.ExchangePerilDirect, fmt.Sprintf("%s.%s", routing.PauseKey, username), routing.PauseKey, pubsub.SimpleQueueTransient)
	if err != nil {
		log.Fatalf("Error subscribing to pause: %v", err)
	}
	log.Printf("Queue %v declared and bound!\n", queue.Name)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)

	<-sigs
	log.Println("Shutting down...")
}
