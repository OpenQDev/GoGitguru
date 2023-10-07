package server

import (
	"main/internal/database"
)

func GetApiConfig(database *database.Queries) (ApiConfig, error) {
	apiCfg := ApiConfig{
		DB:                   database,
		GithubRestAPIBaseUrl: "https://api.github.com",
		GithubGraphQLBaseUrl: "https://api.github.com/graphql",
	}

	return apiCfg, nil
}
