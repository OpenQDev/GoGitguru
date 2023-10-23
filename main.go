package main

import (
	"main/internal/pkg/server"
	"main/internal/pkg/setup"
	"time"
	"usersync"
	"util/logger"
)

func main() {

	env := setup.ExtractAndVerifyEnvironment(".env")

	database, apiCfg := server.PrepareServerSingleton(env.DbUrl)
	logger.SetDebugMode(env.Debug)

	go usersync.StartSyncingUser(database, "repos", 10, time.Duration(env.SyncUsersIntervalMinutesInt)*time.Minute, env.GhAccessToken, 2, apiCfg)

}
