package main

import (
	"main/internal/pkg/handlers"
	"main/internal/pkg/logger"
	"main/internal/pkg/setup"
	"main/internal/pkg/sync"
	"time"
)

func main() {
	portString, dbUrl, originUrl, debugMode := setup.ExtractAndVerifyEnvironment(".env")
	database, apiCfg := setup.PrepareDatabase(dbUrl)
	logger.SetDebugMode(debugMode)
	go sync.StartSyncing(database, "repos", 10, 10*time.Second)
	handlers.StartServer(apiCfg, portString, originUrl)
}
