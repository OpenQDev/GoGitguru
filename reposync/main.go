package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/OpenQDev/GoGitguru/database"
	reposync "github.com/OpenQDev/GoGitguru/reposync/src"
	"github.com/OpenQDev/GoGitguru/util/logger"
	"github.com/OpenQDev/GoGitguru/util/setup"
)

// Message represents the JSON structure of messages consumed from Kafka
type Message struct {
	RepoURL string `json:"repo_url"`
}

type Consumer struct {
	ready    chan bool
	database *database.Queries // Replace with your actual database type
	conn     *sql.DB           // Replace with your actual connection type
	env      setup.EnvConfig   // Replace with your actual connection type
}

func main() {
	env := setup.ExtractAndVerifyEnvironment("../.env")

	database, conn, err := setup.GetDatbase(env.DbUrl)
	if err != nil {
		logger.LogError("Failed to connect to database:", err)
		return
	}
	defer conn.Close()

	logger.SetDebugMode(env.Debug)
	logger.LogBlue("beginning repo syncing...")

	stopChan := make(chan struct{})
	setupSignalHandler(stopChan)

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

	// Spawn consumers
	numConsumers := 5 // Set the number of concurrent consumer instances
	for i := 0; i < numConsumers; i++ {
		wg.Add(1)
		go func(consumerID int) {
			defer wg.Done()
			consumer := Consumer{
				ready:    make(chan bool),
				database: database,
				conn:     conn,
			}

			client, err := sarama.NewConsumerGroup(brokers, group, config)
			if err != nil {
				logger.LogError(fmt.Sprintf("Error creating consumer group client for consumer %d", consumerID), err)
				return
			}
			defer func() {
				if err := client.Close(); err != nil {
					logger.LogError(fmt.Sprintf("Error closing client for consumer %d", consumerID), err)
				}
			}()

			ctx := context.Background()

			for {
				if err := client.Consume(ctx, topics, &consumer); err != nil {
					logger.LogError(fmt.Sprintf("Error from consumer %d", consumerID), err)
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
	<-stopChan
	log.Println("Terminating: via signal")
	wg.Wait()
	log.Println("All consumers have been shut down.")
}

func setupSignalHandler(stopChan chan<- struct{}) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)

	go func() {
		<-sigChan
		close(stopChan)
	}()
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
		var msg Message
		// Parse the JSON-encoded message
		if err := json.Unmarshal(message.Value, &msg); err != nil {
			logger.LogError("Error decoding JSON message", err)
			continue // Skip malformed messages
		}

		// Call StartSyncingCommits with the extracted repo_url
		reposync.StartSyncingCommits(consumer.database, consumer.conn, "repos", consumer.env.GitguruUrl, msg.RepoURL)

		// Mark the message as processed
		session.MarkMessage(message, "")
	}

	// Manually commit offsets after processing messages
	session.Commit()

	return nil
}
