package server

import (
	"database/sql"

	"github.com/OpenQDev/GoGitguru/util/logger"
	"github.com/OpenQDev/GoGitguru/util/setup"
)

func PrepareServerSingleton(dbUrl string) (*sql.DB, ApiConfig) {
	database, conn, err := setup.GetDatbase(dbUrl)

	if err != nil {
		logger.LogError("error getting database: %s", err)
	}

	apiCfg, err := GetApiConfig(database, dbUrl, conn)
	if err != nil {
		logger.LogFatalRedAndExit("can't connect to DB: %s", err)
	}
	return conn, apiCfg
}
