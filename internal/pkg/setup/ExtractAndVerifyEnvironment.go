package setup

import (
	"main/internal/pkg/logger"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func ExtractAndVerifyEnvironment(pathToDotenv string) (portString string, dbUrl string, originUrl string, debug bool, sync bool, syncIntervalMinutesInt int) {
	godotenv.Load(pathToDotenv)
	portString = os.Getenv("PORT")
	dbUrl = os.Getenv("DB_URL")
	originUrl = os.Getenv("ORIGIN_URL")
	debugMode := os.Getenv("DEBUG_MODE")
	syncMode := os.Getenv("SYNC_MODE")
	syncIntervalMinutes := os.Getenv("SYNC_INTERVAL_MINUTES")

	if portString == "" || dbUrl == "" || originUrl == "" || debugMode == "" {
		logger.LogFatalRedAndExit("PORT | DB_URL | ORIGIN_URL | DEBUG_MODE is not found in the environment")
	}

	debug, err := strconv.ParseBool(debugMode)
	if err != nil {
		logger.LogFatalRedAndExit("DEBUG_MODE must be a boolean")
	}

	sync, err = strconv.ParseBool(syncMode)
	if err != nil {
		logger.LogFatalRedAndExit("SYNC_MODE must be a boolean")
	}

	syncIntervalMinutesInt, err = strconv.Atoi(syncIntervalMinutes)
	if err != nil {
		logger.LogFatalRedAndExit("SYNC_INTERVAL_MINUTES must be an integer")
	}

	return portString, dbUrl, originUrl, debug, sync, syncIntervalMinutesInt
}
