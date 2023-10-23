package main

import (
	"time"
	usersync "usersync/src"
	"util/logger"
	"util/setup"
)

func main() {
	env := setup.ExtractAndVerifyEnvironment(".env")

	database, _ := setup.GetDatbase(env.DbUrl)

	logger.SetDebugMode(env.Debug)

	go usersync.StartSyncingUser(database, "repos", 10, time.Duration(env.SyncUsersIntervalMinutesInt)*time.Minute, env.GhAccessToken, 2, "https://github.com/graphql")
}
