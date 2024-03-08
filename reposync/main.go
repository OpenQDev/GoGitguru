package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	reposync "github.com/OpenQDev/GoGitguru/reposync/src"
	"github.com/OpenQDev/GoGitguru/util/logger"
	"github.com/OpenQDev/GoGitguru/util/setup"
)

func main() {
	env := setup.ExtractAndVerifyEnvironment("../.env")

	database, conn, err := setup.GetDatbase(env.DbUrl)
	if err != nil {
		logger.LogError("Failed to connect to database:", err)
		return
	}
	defer conn.Close()

	logger.SetDebugMode(env.Debug)
	logger.LogBlue("beginning repo syncing...")

	stopChan := make(chan struct{})
	setupSignalHandler(stopChan)

	const MAX_CONCURRENT_INSTANCES = 5
	sem := make(chan bool, MAX_CONCURRENT_INSTANCES)

	for {
		select {
		case <-stopChan:
			logger.LogBlue("shutting down gracefully...")
			os.Exit(0)
		default:
			for i := 0; i < MAX_CONCURRENT_INSTANCES; i++ {
				sem <- true
				go func() {
					defer func() { <-sem }()
					reposync.StartSyncingCommits(database, conn, "repos", env.GitguruUrl)
				}()
			}
			time.Sleep(time.Duration(env.RepoSyncInterval) * time.Second)
		}
	}
}

func setupSignalHandler(stopChan chan<- struct{}) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)

	go func() {
		<-sigChan
		close(stopChan)
	}()
}
