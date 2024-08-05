package server

import (
	"database/sql"

	"github.com/IBM/sarama"
	"github.com/OpenQDev/GoGitguru/database"
)

func GetApiConfig(database *database.Queries, dbUrl string, conn *sql.DB) (ApiConfig, error) {
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.Return.Successes = true

	apiCfg := ApiConfig{
		DB:                   database,
		GithubRestAPIBaseUrl: "https://api.github.com",
		GithubGraphQLBaseUrl: "https://api.github.com/graphql",
		PrefixPath:           "./repos",
		DBURL:                dbUrl,
		Conn:                 conn,
		KafkaBrokers:         []string{"localhost:9092"},
		KafkaConfig:          kafkaConfig,
	}

	return apiCfg, nil
}
