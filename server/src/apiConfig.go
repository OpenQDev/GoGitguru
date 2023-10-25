package server

import "github.com/OpenQDev/GoGitguru/database"

type ApiConfig struct {
	DB                   *database.Queries
	GithubRestAPIBaseUrl string
	GithubGraphQLBaseUrl string
	PrefixPath           string
}
