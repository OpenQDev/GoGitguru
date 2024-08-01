package main

import (
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/OpenQDev/GoGitguru/database"
	usersync "github.com/OpenQDev/GoGitguru/usersync/src"
	"github.com/OpenQDev/GoGitguru/util/logger"
	"github.com/OpenQDev/GoGitguru/util/setup"
)

func main() {
	env := setup.ExtractAndVerifyEnvironment("../.env")

	database, conn, err := setup.GetDatbase(env.DbUrl)
	if err != nil {
		logger.LogFatalRedAndExit("unable to connect to database: %s", err)
	}
	defer conn.Close()

	logger.SetDebugMode(env.Debug)

	tokens := strings.Split(env.GhAccessTokens, ",")
	randomToken := tokens[rand.Intn(len(tokens))]

	stopChan := make(chan struct{})
	setupSignalHandler(stopChan)

	go syncUserDependencies(database, env.UserDependenciesSyncInterval, stopChan)
	go syncUsers(database, env.UserSyncInterval, randomToken, stopChan)

	<-stopChan
	logger.LogBlue("shutting down gracefully...")
}

func syncUserDependencies(database *database.Queries, interval int, stopChan <-chan struct{}) {
	for {
		select {
		case <-stopChan:
			return
		default:
			logger.LogBlue("syncing user dependencies...")
			usersync.SyncUserDependencies(database)
			logger.LogBlue("user dependencies synced!")

			time.Sleep(time.Duration(interval) * time.Second)
		}
	}
}

func syncUsers(database *database.Queries, interval int, token string, stopChan <-chan struct{}) {
	for {
		select {
		case <-stopChan:
			return
		default:
			logger.LogBlue("beginning user syncing...")
			logger.LogBlue("syncing commits...")
			usersync.StartUserSyncing(database, "repos", token, 10, "https://api.github.com/graphql")
			logger.LogBlue("commits synced!")

			time.Sleep(time.Duration(interval) * time.Second)
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
