package main

import (
	usersync "github.com/OpenQDev/GoGitguru/usersync/src"
	"github.com/OpenQDev/GoGitguru/util/logger"
	"github.com/OpenQDev/GoGitguru/util/setup"
)

func main() {
	env := setup.ExtractAndVerifyEnvironment(".env")

	database, _ := setup.GetDatbase(env.DbUrl)

	logger.SetDebugMode(env.Debug)

	usersync.StartSyncingUser(database, "repos", env.GhAccessToken, 2, "https://github.com/graphql")
	logger.LogBlue("user sync completed!")
}
