package main

import (
	"fmt"
	"log"

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

	log.Println("Successfully connected to RabbitMQ server...")

	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("Error opening channel: %v", err)
	}

	gamelogic.PrintServerHelp()
mainloop:
	for {
		inputs := gamelogic.GetInput()
		if len(inputs) == 0 {
			continue
		}
		cmd := inputs[0]

		switch cmd {
		case "pause":
			fmt.Println("Publishing paused game state")
			err = pubsub.PublishJSON(
				channel,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				routing.PlayingState{
					IsPaused: true,
				})
			if err != nil {
				log.Fatalf("Error publishing message: %v", err)
			}
		case "resume":
			fmt.Println("Publishing resume game state")
			err = pubsub.PublishJSON(
				channel,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				routing.PlayingState{
					IsPaused: false,
				})
			if err != nil {
				log.Fatalf("Error publishing message: %v", err)
			}
		case "quit":
			log.Println("Shutting down...")
			break mainloop
		default:
			fmt.Println("Invalid command. Please enter another command.")
		}
	}
}
