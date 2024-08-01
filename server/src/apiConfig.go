package server

import (
	"database/sql"

	"github.com/OpenQDev/GoGitguru/database"
)

type ApiConfig struct {
	DB                   *database.Queries
	GithubRestAPIBaseUrl string
	GithubGraphQLBaseUrl string
	PrefixPath           string
	DBURL                string
	Conn                 *sql.DB
}
