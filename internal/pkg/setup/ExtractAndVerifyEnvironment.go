package setup

import (
	"main/internal/pkg/logger"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func ExtractAndVerifyEnvironment(pathToDotenv string) (portString string, dbUrl string, originUrl string, debug bool) {
	godotenv.Load(pathToDotenv)
	portString = os.Getenv("PORT")
	dbUrl = os.Getenv("DB_URL")
	originUrl = os.Getenv("ORIGIN_URL")
	debugMode := os.Getenv("DEBUG_MODE")

	if portString == "" || dbUrl == "" || originUrl == "" || debugMode == "" {
		logger.LogFatalRedAndExit("PORT | DB_URL | ORIGIN_URL | DEBUG_MODE is not found in the environment")
	}

	debug, err := strconv.ParseBool(debugMode)
	if err != nil {
		logger.LogFatalRedAndExit("DEBUG_MODE must be a boolean")
	}

	return portString, dbUrl, originUrl, debug
}
