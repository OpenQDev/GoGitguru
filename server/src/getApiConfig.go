package server

import (
	"database/sql"

	"github.com/OpenQDev/GoGitguru/database"
)

func GetApiConfig(database *database.Queries, dbUrl string, conn *sql.DB) (ApiConfig, error) {
	apiCfg := ApiConfig{
		DB:                   database,
		GithubRestAPIBaseUrl: "https://api.github.com",
		GithubGraphQLBaseUrl: "https://api.github.com/graphql",
		PrefixPath:           "./repos",
		DBURL:                dbUrl,
		Conn:                 conn,
	}

	return apiCfg, nil
}
