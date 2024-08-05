package server

import (
	"crypto/tls"
	"database/sql"

	"github.com/IBM/sarama"
	"github.com/OpenQDev/GoGitguru/database"
)

// setupProducer will create a AsyncProducer and returns it
func setupProducer(environment string, kafkaBrokers []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()

	if environment != "LOCAL" {
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

func GetApiConfig(database *database.Queries, dbUrl string, conn *sql.DB, environment string) (ApiConfig, error) {
	producer, err := setupProducer(environment, []string{"localhost:9092"})
	if err != nil {
		return ApiConfig{}, err
	}

	apiCfg := ApiConfig{
		DB:                   database,
		GithubRestAPIBaseUrl: "https://api.github.com",
		GithubGraphQLBaseUrl: "https://api.github.com/graphql",
		PrefixPath:           "./repos",
		DBURL:                dbUrl,
		Conn:                 conn,
		KafkaBrokers:         []string{"localhost:9092"},
		KafkaProducer:        producer,
	}

	return apiCfg, nil
}
