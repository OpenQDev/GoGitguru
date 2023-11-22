package server

import (
	"github.com/OpenQDev/GoGitguru/database"
)

func GetApiConfig(database *database.Queries, dbUrl string) (ApiConfig, error) {
	apiCfg := ApiConfig{
		DB:                   database,
		GithubRestAPIBaseUrl: "https://api.github.com",
		GithubGraphQLBaseUrl: "https://api.github.com/graphql",
		PrefixPath:           "./repos",
		DBURL:                dbUrl,
	}

	return apiCfg, nil
}
