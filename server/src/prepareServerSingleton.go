package server

import (
	"database/database"
	"util/logger"
	"util/setup"
)

func PrepareServerSingleton(dbUrl string) (*database.Queries, ApiConfig) {
	database, err := setup.GetDatbase(dbUrl)

	if err != nil {
		logger.LogError("error getting database: %s", err)
	}

	apiCfg, err := GetApiConfig(database)
	if err != nil {
		logger.LogFatalRedAndExit("can't connect to DB: %s", err)
	}
	return database, apiCfg
}
