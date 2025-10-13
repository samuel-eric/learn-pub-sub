package main

import (
	"log"
	"os"
	"os/signal"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	connStr := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connStr)
	if err != nil {
		log.Fatalf("Error connecting to RabbitMQ server: %v", err)
	}
	defer conn.Close()

	log.Println("Successfully connected to RabbitMQ server...")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)

	<-sigs
	log.Println("Shutting down...")
}
