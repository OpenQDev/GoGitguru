package main

import (
	"main/internal/pkg/logger"
	"main/internal/pkg/server"
	"main/internal/pkg/setup"
	"main/internal/pkg/sync"
	"time"
)

func main() {
	portString, dbUrl, originUrl, debugMode, syncMode, syncIntervalMinutes := setup.ExtractAndVerifyEnvironment(".env")
	database, apiCfg := setup.PrepareDatabase(dbUrl)
	logger.SetDebugMode(debugMode)
	if syncMode {
		go sync.StartSyncing(database, "repos", 10, time.Duration(syncIntervalMinutes)*time.Minute)
	}
	server.StartServer(apiCfg, portString, originUrl)
}
