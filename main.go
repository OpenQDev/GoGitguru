package main

import (
	"main/internal/pkg/logger"
	"main/internal/pkg/server"
	"main/internal/pkg/setup"
	"main/internal/pkg/sync"
	"time"
)

func main() {

	portString,
		dbUrl,
		originUrl,
		debugMode,
		syncMode,
		syncIntervalMinutes,
		syncUsersMode,
		syncUsersIntervalMinutes, ghAccessToken := setup.ExtractAndVerifyEnvironment(".env")

	database, apiCfg := setup.PrepareDatabase(dbUrl)
	logger.SetDebugMode(debugMode)

	if syncMode {
		go sync.StartSyncingCommits(database, "repos", 10, time.Duration(syncIntervalMinutes)*time.Minute)
	}

	if syncUsersMode {
		go sync.StartSyncingUser(database, "repos", 10, time.Duration(syncUsersIntervalMinutes)*time.Minute, ghAccessToken)
	}

	server.StartServer(apiCfg, portString, originUrl)
}
