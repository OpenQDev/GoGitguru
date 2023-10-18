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

	envConfig := setup.ExtractAndVerifyEnvironment(".env")

	database, apiCfg := server.PrepareServerSingleton(envConfig.DbUrl)
	logger.SetDebugMode(envConfig.Debug)

	if envConfig.Sync {
		go reposync.StartSyncingCommits(database, "repos", 10, time.Duration(envConfig.SyncIntervalMinutesInt)*time.Minute)
	}

	if envConfig.SyncUsers {
		time.Sleep(3 * time.Second)
		go usersync.StartSyncingUser(database, "repos", 10, time.Duration(envConfig.SyncUsersIntervalMinutesInt)*time.Minute, envConfig.GhAccessToken, 2, apiCfg)
	}

	if envConfig.StartServer {
		server.StartServer(apiCfg, envConfig.PortString, envConfig.OriginUrl)
	}
}
