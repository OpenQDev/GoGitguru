package main

import (
	"main/internal/pkg/logger"
	"main/internal/pkg/setup"
	"time"
)

func main() {
	portString, dbUrl, originUrl, debugMode := setup.ExtractAndVerifyEnvironment(".env")
	database, apiCfg := setup.PrepareDatabase(dbUrl)
	logger.SetDebugMode(debugMode)
	go setup.StartSyncing(database, "repos", 10, 10*time.Second)
	setup.StartServer(apiCfg, portString, originUrl)
}
