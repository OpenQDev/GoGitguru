package main

import (
	reposync "reposync/src"
	"time"
	"util/logger"
	"util/setup"
)

func main() {
	env := setup.ExtractAndVerifyEnvironment(".env")

	database, _ := setup.GetDatbase(env.DbUrl)

	logger.SetDebugMode(env.Debug)

	reposync.StartSyncingCommits(database, "repos", 10, time.Duration(env.SyncIntervalMinutesInt)*time.Minute)
}
