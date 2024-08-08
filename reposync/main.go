package main

import (
	"context"
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/IBM/sarama"
	"github.com/OpenQDev/GoGitguru/database"
	reposync "github.com/OpenQDev/GoGitguru/reposync/src"
	kafkahelpers "github.com/OpenQDev/GoGitguru/util/kafkaHelpers"
	"github.com/OpenQDev/GoGitguru/util/logger"
	"github.com/OpenQDev/GoGitguru/util/setup"
)

// Message represents the JSON structure of messages consumed from Kafka
type Message struct {
	RepoURL string `json:"repo_url"`
}

type Consumer struct {
	ready    chan bool
	database *database.Queries
	conn     *sql.DB
	env      setup.EnvConfig
	producer sarama.SyncProducer
}

// setupProducer will create a SyncProducer and returns it
func setupProducer(environment string, kafkaBrokers []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	fmt.Printf("Starting producer for environment %s\n", environment)
	if environment == "production" {
		// Set the SASL/OAUTHBEARER configuration
		config.Net.SASL.Enable = true
		config.Net.SASL.Mechanism = sarama.SASLTypeOAuth
		config.Net.SASL.TokenProvider = &reposync.MSKAccessTokenProvider{}

		// Enable TLS
		tlsConfig := tls.Config{}
		config.Net.TLS.Enable = true
		config.Net.TLS.Config = &tlsConfig
	}

	return sarama.NewSyncProducer(kafkaBrokers, config)
}

func main() {
	env := setup.ExtractAndVerifyEnvironment("../.env")

	database, conn, err := setup.GetDatbase(env.DbUrl)
	if err != nil {
		logger.LogError("Failed to connect to database:", err)
		return
	}
	defer conn.Close()

	group := env.RepoUrlsConsumerGroup
	brokers := strings.Split(env.KafkaBrokerUrls, ",")
	topics := []string{env.RepoUrlsTopic}
	producer, err := kafkahelpers.SetupProducer(env.Environment, brokers)
	if err != nil {
		logger.LogError("Failed to setup Kafka producer:", err)
	}

	logger.SetDebugMode(env.Debug)
	logger.LogBlue("beginning repo syncing...")

	stopChan := make(chan struct{})
	setup.SignalHandler(stopChan)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a wait group to manage goroutines
	var wg sync.WaitGroup

	// Spawn consumers
	numConsumers := env.RepoSyncConsumerCount
	for i := 0; i < numConsumers; i++ {
		wg.Add(1)
		go func(consumerID int) {
			defer wg.Done()
			consumer := Consumer{
				ready:    make(chan bool),
				database: database,
				conn:     conn,
				env:      env,
				producer: producer,
			}

			client, err := kafkahelpers.SetUpConsumerGroup(env.Environment, brokers, group)
			if err != nil {
				logger.LogError(fmt.Sprintf("Error creating consumer group client for consumer %d", consumerID), err)
				return
			}
			defer func() {
				if err := client.Close(); err != nil {
					logger.LogError(fmt.Sprintf("Error closing client for consumer %d", consumerID), err)
				}
			}()

			for {
				select {
				case <-ctx.Done():
					return
				default:
					if err := client.Consume(ctx, topics, &consumer); err != nil {
						logger.LogError(fmt.Sprintf("Error from consumer %d", consumerID), err)
					}
				}
				consumer.ready = make(chan bool)
			}
		}(i)
	}

	log.Println("All consumers are up and running!")
	<-stopChan
	log.Println("Terminating: via signal")
	cancel() // Cancel the context to stop all consumers
	wg.Wait()
	log.Println("All consumers have been shut down.")
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	close(consumer.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	consumer.producer.Close()
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				return nil // Channel closed
			}
			var msg Message
			if err := json.Unmarshal(message.Value, &msg); err != nil {
				logger.LogError("Error decoding JSON message", err)
				continue
			}

			reposync.StartSyncingCommits(consumer.database, consumer.conn, "repos", consumer.env.GitguruUrl, msg.RepoURL, consumer.producer)

			session.MarkMessage(message, "")
			session.Commit()

		case <-session.Context().Done():
			return nil // Session canceled
		}
	}
}
