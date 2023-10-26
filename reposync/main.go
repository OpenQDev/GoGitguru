package main

import (
	"time"

	reposync "github.com/OpenQDev/GoGitguru/reposync/src"
	"github.com/OpenQDev/GoGitguru/util/logger"
	"github.com/OpenQDev/GoGitguru/util/setup"
)

func main() {
	env := setup.ExtractAndVerifyEnvironment(".env")

	database, _ := setup.GetDatbase(env.DbUrl)

	logger.SetDebugMode(env.Debug)

	for {
		logger.LogBlue("beginning repo syncing...")
		reposync.StartSyncingCommits(database, "repos", 10, time.Duration(env.SyncIntervalMinutesInt)*time.Minute)
		time.Sleep(time.Minute)
		logger.LogBlue("repo sync completed!")
	}
}
