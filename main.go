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

	env := setup.ExtractAndVerifyEnvironment(".env")

	database, apiCfg := server.PrepareServerSingleton(env.DbUrl)
	logger.SetDebugMode(env.Debug)

	if env.Sync {
		go reposync.StartSyncingCommits(database, "repos", 10, time.Duration(env.SyncIntervalMinutesInt)*time.Minute)
	}

	if env.SyncUsers {
		time.Sleep(3 * time.Second)
		go usersync.StartSyncingUser(database, "repos", 10, time.Duration(env.SyncUsersIntervalMinutesInt)*time.Minute, env.GhAccessToken, 2, apiCfg)
	}

	if env.StartServer {
		server.StartServer(apiCfg, env.PortString, env.OriginUrl)
	}
}
