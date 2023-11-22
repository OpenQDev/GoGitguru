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

	logger.LogBlue("beginning repo syncing...")

	// PRODUCTION: This runs as a CronJob on Kubernetes. Therefore, it's interval is set by the CRON_STRING parameter
	// DEVELOPMENT: To mimic the interval, here we check for the REPOSYNC_INTERVAL environment variable to periodically re-run StartSyncingCommits

	if env.RepoSyncInterval != 0 {
		for {
			for i := 0; i < 2; i++ {
				go reposync.StartSyncingCommits(database, "repos", env.GitguruUrl, env.GetDueRepoUrlExpiration)
			}
			time.Sleep(time.Duration(env.RepoSyncInterval) * time.Second)
		}
	} else {
		reposync.StartSyncingCommits(database, "repos", env.GitguruUrl, env.GetDueRepoUrlExpiration)
	}

	logger.LogBlue("repo sync completed!")
}
