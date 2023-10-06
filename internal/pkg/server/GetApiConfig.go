package server

import (
	"main/internal/database"
)

func GetApiConfig(database *database.Queries) (ApiConfig, error) {
	apiCfg := ApiConfig{
		DB: database,
	}

	return apiCfg, nil
}
