package main

import (
	"math/rand"
	"strings"
	"time"

	usersync "github.com/OpenQDev/GoGitguru/usersync/src"
	"github.com/OpenQDev/GoGitguru/util/logger"
	"github.com/OpenQDev/GoGitguru/util/setup"
)

func main() {
	env := setup.ExtractAndVerifyEnvironment(".env")

	database, conn, err := setup.GetDatbase(env.DbUrl)
	if err != nil {
		logger.LogFatalRedAndExit("unable to connect to database: %s", err)
	}
	defer conn.Close()

	logger.SetDebugMode(env.Debug)

	logger.LogBlue("beginning user syncing...")

	// PRODUCTION: This runs as a CronJob on Kubernetes. Therefore, it's interval is set by the CRON_STRING parameter
	// DEVELOPMENT: To mimic the interval, here we check for the USERSYNC_INTERVAL environment variable to periodically re-run StartSyncingUser

	tokens := strings.Split(env.GhAccessTokens, ",")
	randomToken := tokens[rand.Intn(len(tokens))]

	if env.UserSyncInterval != 0 {
		for {
			usersync.StartSyncingUser(database, "repos", randomToken, 10, "https://api.github.com/graphql")
			time.Sleep(time.Duration(env.UserSyncInterval) * time.Second)
		}
	} else {
		usersync.StartSyncingUser(database, "repos", randomToken, 10, "https://api.github.com/graphql")
	}

	logger.LogBlue("user sync completed!")
}
