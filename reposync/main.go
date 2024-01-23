package main

import (
	"time"

	reposync "github.com/OpenQDev/GoGitguru/reposync/src"
	"github.com/OpenQDev/GoGitguru/util/logger"
	"github.com/OpenQDev/GoGitguru/util/setup"
)

func main() {
	env := setup.ExtractAndVerifyEnvironment("../.env")

	database, conn, _ := setup.GetDatbase(env.DbUrl)
	defer conn.Close()

	logger.SetDebugMode(env.Debug)

	logger.LogBlue("beginning repo syncing...")

	// PRODUCTION: This runs as a CronJob on Kubernetes. Therefore, it's interval is set by the CRON_STRING parameter
	// DEVELOPMENT: To mimic the interval, here we check for the REPOSYNC_INTERVAL environment variable to periodically re-run StartSyncingCommits

	const MAX_CONCURRENT_INSTANCES = 5
	if env.RepoSyncInterval != 0 {
		sem := make(chan bool, MAX_CONCURRENT_INSTANCES) // create a buffered channel with capacity MAX_CONCURRENT_INSTANCES
		for {
			for i := 0; i < MAX_CONCURRENT_INSTANCES; i++ {
				sem <- true // block if there are already MAX_CONCURRENT_INSTANCES goroutines running
				go func() {
					reposync.StartSyncingCommits(database, conn, "repos", env.GitguruUrl)
					<-sem // release the semaphore when goroutine finishes
				}()
			}
			time.Sleep(time.Duration(env.RepoSyncInterval) * time.Second)
		}
	} else {
		reposync.StartSyncingCommits(database, conn, "repos", env.GitguruUrl)
	}

	logger.LogBlue("repo sync completed!")
}
