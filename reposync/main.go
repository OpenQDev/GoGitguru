package main

import (
	reposync "github.com/OpenQDev/GoGitguru/reposync/src"
	"github.com/OpenQDev/GoGitguru/util/logger"
	"github.com/OpenQDev/GoGitguru/util/setup"
)

func main() {
	env := setup.ExtractAndVerifyEnvironment(".env")

	database, _ := setup.GetDatbase(env.DbUrl)

	logger.SetDebugMode(env.Debug)

	logger.LogBlue("beginning repo syncing...")
	reposync.StartSyncingCommits(database, "repos")
	logger.LogBlue("repo sync completed!")
}
