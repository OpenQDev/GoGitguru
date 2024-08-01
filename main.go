package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/IBM/sarama"
)

type Consumer struct {
	ready chan bool
}

func main() {
	// Set up the Sarama configuration
	config := sarama.NewConfig()
	config.Version = sarama.V2_5_0_0 // Set the version to match your Kafka cluster
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRange()
	config.Consumer.Offsets.Initial = sarama.OffsetNewest // Start from the newest message
	config.Consumer.Offsets.AutoCommit.Enable = false     // Disable auto commit for manual control

	// Define the consumer group and brokers
	group := "repo-urls-group"
	brokers := []string{"localhost:9092"} // Replace with your broker addresses
	topics := []string{"repo-urls"}

	// Create a wait group to manage goroutines
	var wg sync.WaitGroup

	// Set up signal handling for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	// Spawn 5 consumers
	numConsumers := 5
	for i := 0; i < numConsumers; i++ {
		wg.Add(1)
		go func(consumerID int) {
			defer wg.Done()
			consumer := Consumer{
				ready: make(chan bool),
			}

			client, err := sarama.NewConsumerGroup(brokers, group, config)
			if err != nil {
				log.Fatalf("Error creating consumer group client for consumer %d: %v\n", consumerID, err)
			}
			defer func() {
				if err := client.Close(); err != nil {
					log.Printf("Error closing client for consumer %d: %v\n", consumerID, err)
				}
			}()

			for {
				if err := client.Consume(ctx, topics, &consumer); err != nil {
					log.Printf("Error from consumer %d: %v\n", consumerID, err)
				}
				// Check if context was canceled, signaling that the consumer should stop
				if ctx.Err() != nil {
					return
				}
				consumer.ready = make(chan bool)
			}
		}(i)
	}

	log.Println("All consumers are up and running!")
	<-sigterm
	log.Println("Terminating: via signal")
	cancel()
	wg.Wait()
	log.Println("All consumers have been shut down.")
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		fmt.Printf("Consumer %s received message: %s\n", session.MemberID(), string(message.Value))
		// Mark the message as processed
		session.MarkMessage(message, "")
	}

	// Manually commit offsets after processing messages
	session.Commit()

	return nil
}
