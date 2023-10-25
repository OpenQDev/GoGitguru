package main

import (
	"time"

	usersync "github.com/OpenQDev/GoGitguru/usersync/src"
	"github.com/OpenQDev/GoGitguru/util/logger"
	"github.com/OpenQDev/GoGitguru/util/setup"
)

func main() {
	env := setup.ExtractAndVerifyEnvironment(".env")

	database, _ := setup.GetDatbase(env.DbUrl)

	logger.SetDebugMode(env.Debug)

	for {
		usersync.StartSyncingUser(database, "repos", 10, time.Duration(env.SyncUsersIntervalMinutesInt)*time.Minute, env.GhAccessToken, 2, "https://github.com/graphql")
		time.Sleep(time.Minute)
	}
}
