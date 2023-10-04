package setup

import (
	"main/internal/database"
	"main/internal/pkg/handlers"
)

func GetApiConfig(database *database.Queries) (handlers.ApiConfig, error) {
	apiCfg := handlers.ApiConfig{
		DB: database,
	}

	return apiCfg, nil
}
