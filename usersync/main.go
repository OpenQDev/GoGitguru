package main

import (
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

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

	for {
		select {
		case <-stopChan:
			logger.LogBlue("shutting down gracefully...")
			os.Exit(0)
		default:
			logger.LogBlue("beginning user syncing...")
			usersync.StartSyncingUser(database, "repos", randomToken, 10, "https://api.github.com/graphql")
			time.Sleep(time.Duration(env.UserSyncInterval) * time.Second)
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
