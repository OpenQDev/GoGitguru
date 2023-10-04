package setup

import (
	"main/internal/database"
	"main/internal/pkg/handlers"
	"main/internal/pkg/logger"
)

func PrepareDatabase(dbUrl string) (*database.Queries, handlers.ApiConfig) {
	database, err := GetDatbase(dbUrl)
	if err != nil {
		logger.LogError("error getting database: %s", err)
	}

	apiCfg, err := GetApiConfig(database)
	if err != nil {
		logger.LogFatalRedAndExit("can't connect to DB: %s", err)
	}
	return database, apiCfg
}
