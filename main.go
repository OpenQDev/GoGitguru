package main

import (
	"main/internal/pkg/logger"
	"main/internal/pkg/reposync"
	"main/internal/pkg/server"
	"main/internal/pkg/setup"
	"main/internal/pkg/usersync"
	"time"
)

func main() {

	portString,
		dbUrl,
		originUrl,
		debugMode,
		repoSyncMode,
		syncIntervalMinutes,
		syncUsersMode,
		syncUsersIntervalMinutes,
		ghAccessToken,
		_,
		startServer := setup.ExtractAndVerifyEnvironment(".env")

	database, apiCfg := server.PrepareServerSingleton(dbUrl)
	logger.SetDebugMode(debugMode)

	if repoSyncMode {
		go reposync.StartSyncingCommits(database, "repos", 10, time.Duration(syncIntervalMinutes)*time.Minute)
	}

	if syncUsersMode {
		time.Sleep(3 * time.Second)
		go usersync.StartSyncingUser(database, "repos", 10, time.Duration(syncUsersIntervalMinutes)*time.Minute, ghAccessToken, 2, apiCfg)
	}

	if startServer {
		server.StartServer(apiCfg, portString, originUrl)
	}
}
