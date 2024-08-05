package server

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"strings"

	"github.com/IBM/sarama"
	"github.com/OpenQDev/GoGitguru/database"
	"github.com/OpenQDev/GoGitguru/util/logger"
	"github.com/OpenQDev/GoGitguru/util/setup"
)

// setupProducer will create a SyncProducer and returns it
func setupProducer(environment string, kafkaBrokers []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	fmt.Printf("Starting producer for environment %s\n", environment)
	if environment == "production" {
		// Set the SASL/OAUTHBEARER configuration
		config.Net.SASL.Enable = true
		config.Net.SASL.Mechanism = sarama.SASLTypeOAuth
		config.Net.SASL.TokenProvider = &MSKAccessTokenProvider{}

		// Enable TLS
		tlsConfig := tls.Config{}
		config.Net.TLS.Enable = true
		config.Net.TLS.Config = &tlsConfig
	}

	return sarama.NewSyncProducer(kafkaBrokers, config)
}

func GetApiConfig(database *database.Queries, dbUrl string, conn *sql.DB, environment setup.EnvConfig) (ApiConfig, error) {
	kafkaBrokers := strings.Split(environment.KafkaBrokerUrls, ",")
	producer, err := setupProducer(environment.Environment, kafkaBrokers)
	if err != nil {
		logger.LogFatalRedAndExit("Failed to setup Kafka producer: %v", err)
	}

	apiCfg := ApiConfig{
		DB:                   database,
		GithubRestAPIBaseUrl: "https://api.github.com",
		GithubGraphQLBaseUrl: "https://api.github.com/graphql",
		PrefixPath:           "./repos",
		DBURL:                dbUrl,
		Conn:                 conn,
		KafkaProducer:        producer,
	}

	return apiCfg, nil
}
