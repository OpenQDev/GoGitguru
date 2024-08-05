package main

import (
	"context"
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/IBM/sarama"

	"github.com/OpenQDev/GoGitguru/database"
	reposync "github.com/OpenQDev/GoGitguru/reposync/src"
	usersync "github.com/OpenQDev/GoGitguru/usersync/src"
	"github.com/OpenQDev/GoGitguru/util/logger"
	"github.com/OpenQDev/GoGitguru/util/setup"
)

type Consumer struct {
	ready    chan bool
	database *database.Queries // Replace with your actual database type
	conn     *sql.DB           // Replace with your actual connection type
	env      setup.EnvConfig   // Replace with your actual connection type
}

func setUpConsumerGroup(environment string, kafkaBrokers []string, group string) (sarama.ConsumerGroup, error) {
	// Set the SASL/OAUTHBEARER configuration
	// Set up the Sarama configuration
	config := sarama.NewConfig()
	config.Version = sarama.V2_5_0_0
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRange()
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Offsets.AutoCommit.Enable = false

	fmt.Printf("Starting consumer for environment %s\n", environment)
	if environment == "production" {
		config.Net.SASL.Enable = true
		config.Net.SASL.Mechanism = sarama.SASLTypeOAuth
		config.Net.SASL.TokenProvider = &reposync.MSKAccessTokenProvider{}

		tlsConfig := tls.Config{}
		config.Net.TLS.Enable = true
		config.Net.TLS.Config = &tlsConfig
	}

	return sarama.NewConsumerGroup(kafkaBrokers, group, config)
}

func main() {
	env := setup.ExtractAndVerifyEnvironment("../.env")

	database, conn, err := setup.GetDatbase(env.DbUrl)
	if err != nil {
		logger.LogFatalRedAndExit("unable to connect to database: %s", err)
	}
	defer conn.Close()

	logger.SetDebugMode(env.Debug)

	stopChan := make(chan struct{})
	setupSignalHandler(stopChan)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Define the consumer group and brokers
	group := env.UserSyncConsumerGroup
	brokers := strings.Split(env.KafkaBrokerUrls, ",")
	topics := []string{env.UserSyncTopic}

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

			client, err := setUpConsumerGroup(env.Environment, brokers, group)
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

	logger.LogBlue("All consumers are up and running!")

	<-stopChan
	logger.LogBlue("Terminating: via signal")
	cancel()
	wg.Wait()
	logger.LogBlue("All consumers have been shut down.")
}

/*
func syncUserDependencies(database *database.Queries, interval int, stopChan <-chan struct{}) {
	for {
		select {
		case <-stopChan:
			return
		default:
			logger.LogBlue("syncing user dependencies...")
			usersync.SyncUserDependencies(database)
			logger.LogBlue("user dependencies synced!")

			time.Sleep(time.Duration(interval) * time.Second)
		}
	}
}

func syncUsers(database *database.Queries, interval int, token string, stopChan <-chan struct{}, msg usersync.Message) {
	for {
		select {
		case <-stopChan:
			return
		default:
			logger.LogBlue("beginning user syncing...")
			logger.LogBlue("syncing commits...")
			usersync.StartUserSyncing(database, "repos", token, 10, "https://api.github.com/graphql", msg)
			logger.LogBlue("commits synced!")

			time.Sleep(time.Duration(interval) * time.Second)
		}
	}
}*/

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

	env := setup.ExtractAndVerifyEnvironment("../.env")

	logger.SetDebugMode(env.Debug)

	tokens := strings.Split(env.GhAccessTokens, ",")
	token := tokens[rand.Intn(len(tokens))]
	for message := range claim.Messages() {
		var msg usersync.Message
		fmt.Println(string(message.Value))
		// Parse the JSON-encoded message
		if err := json.Unmarshal(message.Value, &msg); err != nil {
			logger.LogError("Error decoding JSON message for usersync", err)
			continue // Skip malformed messages
		}

		// Call StartSyncingCommits with the extracted repo_url
		usersync.StartUserSyncing(consumer.database, "repos", token, 10, "https://api.github.com/graphql", msg)

		// Mark the message as processed
		session.MarkMessage(message, "")
	}

	// Manually commit offsets after processing messages
	session.Commit()

	return nil
}
