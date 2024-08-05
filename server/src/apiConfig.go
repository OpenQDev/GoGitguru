package server

import (
	"database/sql"

	"github.com/IBM/sarama"
	"github.com/OpenQDev/GoGitguru/database"
)

type ApiConfig struct {
	DB                   *database.Queries
	GithubRestAPIBaseUrl string
	GithubGraphQLBaseUrl string
	PrefixPath           string
	DBURL                string
	Conn                 *sql.DB
	KafkaBrokers         []string
	KafkaConfig          *sarama.Config
}
