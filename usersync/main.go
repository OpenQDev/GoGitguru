package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"

	"github.com/IBM/sarama"

	"github.com/OpenQDev/GoGitguru/database"
	usersync "github.com/OpenQDev/GoGitguru/usersync/src"
	kafkahelpers "github.com/OpenQDev/GoGitguru/util/kafkaHelpers"
	"github.com/OpenQDev/GoGitguru/util/logger"
	"github.com/OpenQDev/GoGitguru/util/setup"
)

type TopicParams struct {
	Topic         string
	ConsumerGroup string
	ConsumerCount int
}
type Consumer struct {
	ready    chan bool
	database *database.Queries // Replace with your actual database type
	conn     *sql.DB           // Replace with your actual connection type
	env      setup.EnvConfig   // Replace with your actual connection type
	producer sarama.SyncProducer
}

func handleConsumerGroup(wg *sync.WaitGroup, topicParams TopicParams, env setup.EnvConfig) {

	database, conn, err := setup.GetDatbase(env.DbUrl)
	if err != nil {
		logger.LogFatalRedAndExit("unable to connect to database: %s", err)
	}
	defer conn.Close()

	brokers := strings.Split(env.KafkaBrokerUrls, ",")
	producer, err := kafkahelpers.SetupProducer(env.Environment, brokers)
	if err != nil {
		logger.LogError("Failed to setup Kafka producer:", err)
	}

	logger.SetDebugMode(env.Debug)

	stopChan := make(chan struct{})
	setup.SignalHandler(stopChan)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Define the consumer group and brokers
	group := topicParams.ConsumerGroup
	topics := []string{topicParams.Topic}

	// Spawn consumers
	numConsumers := topicParams.ConsumerCount // Set the number of concurrent consumer instances
	for i := 0; i < numConsumers; i++ {
		wg.Add(1)
		go func(consumerID int) {
			defer wg.Done()
			consumer := Consumer{
				ready:    make(chan bool),
				database: database,
				conn:     conn,
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

	logger.LogBlue("All consumers are up and running!")

	<-stopChan
	logger.LogBlue("Terminating: via signal")
	cancel()
	wg.Wait()
	logger.LogBlue("All consumers have been shut down.")
}

func main() {
	var wg sync.WaitGroup

	env := setup.ExtractAndVerifyEnvironment("../.env")
	wg.Add(2)
	go func() {
		topicParams := TopicParams{
			Topic:         env.UserDepsSyncTopic,
			ConsumerGroup: env.UserDepsSyncConsumerGroup,
			ConsumerCount: env.UserDepsSyncConsumerCount,
		}
		handleConsumerGroup(&wg, topicParams, env)
		wg.Done()
		os.Exit(0)
	}()
	go func() {
		topicParams := TopicParams{
			Topic:         env.UserSyncTopic,
			ConsumerGroup: env.UserSyncConsumerGroup,
			ConsumerCount: env.UserSyncConsumerCount,
		}
		handleConsumerGroup(&wg, topicParams, env)
		wg.Done()
		os.Exit(0)
	}()

	wg.Wait()
	os.Exit(0)
	logger.LogBlue("All consumers have been shut down.")
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
		fmt.Println(string(message.Value), "recieved")

		if claim.Topic() == env.UserSyncTopic {
			var msg usersync.Message
			// Parse the JSON-encoded message
			if err := json.Unmarshal(message.Value, &msg); err != nil {
				logger.LogError("Error decoding JSON message for usersync", err)
				continue // Skip malformed messages
			}

			// Call StartSyncingCommits with the extracted repo_url
			usersync.StartUserSyncing(consumer.database, "repos", token, 10, "https://api.github.com/graphql", msg, consumer.producer)

			// Mark the message as processed
			session.MarkMessage(message, "")
		}

		if claim.Topic() == env.UserDepsSyncTopic {
			var msg usersync.DepsMessage

			// Parse the JSON-encoded message
			if err := json.Unmarshal(message.Value, &msg); err != nil {
				logger.LogError("Error decoding JSON message for usersync", err)
				continue // Skip malformed messages
			}

			// Call StartSyncingCommits with the extracted repo_url
			usersync.SyncUserDependencies(consumer.database, msg.RepoUrl)

			// Mark the message as processed
			session.MarkMessage(message, "")

		}
	}

	// Manually commit offsets after processing messages
	session.Commit()

	return nil
}
