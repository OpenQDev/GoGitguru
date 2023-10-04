package setup

import (
	"main/internal/database"
	"main/internal/pkg/server"
)

func GetApiConfig(database *database.Queries) (server.ApiConfig, error) {
	apiCfg := server.ApiConfig{
		DB: database,
	}

	return apiCfg, nil
}
