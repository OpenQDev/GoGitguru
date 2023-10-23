package server

import "database/database"

type ApiConfig struct {
	DB                   *database.Queries
	GithubRestAPIBaseUrl string
	GithubGraphQLBaseUrl string
	PrefixPath           string
}
