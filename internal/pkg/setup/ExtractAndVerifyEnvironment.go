package setup

import (
	"main/internal/pkg/logger"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func ExtractAndVerifyEnvironment(
	pathToDotenv string,
) (
	portString string,
	dbUrl string,
	originUrl string,
	debug bool,
	sync bool,
	syncIntervalMinutesInt int,
	syncUsers bool,
	syncUsersIntervalMinutesInt int,
	ghAccessToken string,
) {
	godotenv.Load(pathToDotenv)
	portString = os.Getenv("PORT")
	dbUrl = os.Getenv("DB_URL")
	originUrl = os.Getenv("ORIGIN_URL")
	debugMode := os.Getenv("DEBUG_MODE")
	syncMode := os.Getenv("SYNC_COMMITS_MODE")
	syncIntervalMinutes := os.Getenv("SYNC_COMMITS_INTERVAL_MINUTES")
	syncUsersMode := os.Getenv("SYNC_USERS_MODE")
	syncUsersIntervalMinutes := os.Getenv("SYNC_USERS_INTERVAL_MINUTES")
	ghAccessToken = os.Getenv("GH_ACCESS_TOKEN")

	if portString == "" || dbUrl == "" || originUrl == "" || debugMode == "" {
		logger.LogFatalRedAndExit("PORT | DB_URL | ORIGIN_URL | DEBUG_MODE is not found in the environment")
	}

	debug, err := strconv.ParseBool(debugMode)
	if err != nil {
		logger.LogFatalRedAndExit("DEBUG_MODE must be a boolean")
	}

	sync, err = strconv.ParseBool(syncMode)
	if err != nil {
		logger.LogFatalRedAndExit("SYNC_COMMITS_MODE must be a boolean")
	}

	syncIntervalMinutesInt, err = strconv.Atoi(syncIntervalMinutes)
	if err != nil {
		logger.LogFatalRedAndExit("SYNC_INTERVAL_MINUTES must be an integer")
	}

	syncUsers, err = strconv.ParseBool(syncUsersMode)
	if err != nil {
		logger.LogFatalRedAndExit("SYNC_USERS_MODE must be a boolean")
	}

	syncUsersIntervalMinutesInt, err = strconv.Atoi(syncUsersIntervalMinutes)
	if err != nil {
		logger.LogFatalRedAndExit("SYNC_INTERVAL_MINUTES must be an integer")
	}

	return portString, dbUrl, originUrl, debug, sync, syncIntervalMinutesInt, syncUsers, syncUsersIntervalMinutesInt, ghAccessToken
}
