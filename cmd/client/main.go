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

	log.Println("Starting Peril client...")

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	gameState := gamelogic.NewGameState(username)

	queueName := fmt.Sprintf("%s.%s", routing.PauseKey, username)
	err = pubsub.SubscribeJSON(conn, routing.ExchangePerilDirect, queueName, routing.PauseKey, pubsub.SimpleQueueTransient, handlerPause(gameState))
	if err != nil {
		log.Fatalf("Error subscribing to pause: %v", err)
	}
	log.Printf("Queue %v declared and bound!\n", queueName)

	queueName = fmt.Sprintf("%s.%s", routing.ArmyMovesPrefix, username)
	err = pubsub.SubscribeJSON(conn, routing.ExchangePerilTopic, queueName, fmt.Sprintf("%s.*", routing.ArmyMovesPrefix), pubsub.SimpleQueueTransient, handlerMove(gameState))
	if err != nil {
		log.Fatalf("Error subscribing to pause: %v", err)
	}
	log.Printf("Queue %v declared and bound!\n", queueName)

	moveChannel, err := conn.Channel()
	if err != nil {
		log.Fatalf("Error creating move channel: %v", err)
	}

	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}

		switch words[0] {
		case "spawn":
			err = gameState.CommandSpawn(words)
			if err != nil {
				fmt.Println(err)
				continue
			}
		case "move":
			am, err := gameState.CommandMove(words)
			if err != nil {
				fmt.Println(err)
				continue
			}
			err = pubsub.PublishJSON(moveChannel, routing.ExchangePerilTopic, fmt.Sprintf("%s.%s", routing.ArmyMovesPrefix, username), am)
			if err != nil {
				fmt.Printf("Error publishing message: %v\n", err)
			}
			fmt.Println("Move published successfully")
		case "status":
			gameState.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			fmt.Println("Spamming not allowed yet!")
		case "quit":
			gamelogic.PrintQuit()
			return
		default:
			fmt.Println("Invalid command. Please enter another command.")
		}
	}
}
