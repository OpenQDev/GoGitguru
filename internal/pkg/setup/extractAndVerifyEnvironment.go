package setup

import (
	"main/internal/pkg/logger"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
	PortString                  string
	DbUrl                       string
	OriginUrl                   string
	Debug                       bool
	Sync                        bool
	SyncIntervalMinutesInt      int
	SyncUsers                   bool
	SyncUsersIntervalMinutesInt int
	GhAccessToken               string
	TargetLiveGithub            bool
	StartServer                 bool
}

func ExtractAndVerifyEnvironment(pathToDotenv string) EnvConfig {
	godotenv.Load(pathToDotenv)

	return EnvConfig{
		PortString:                  getEnvVar("PORT", "string").(string),
		DbUrl:                       getEnvVar("DB_URL", "string").(string),
		OriginUrl:                   getEnvVar("ORIGIN_URL", "string").(string),
		Debug:                       getEnvVar("DEBUG_MODE", "bool").(bool),
		Sync:                        getEnvVar("SYNC_COMMITS_MODE", "bool").(bool),
		SyncIntervalMinutesInt:      getEnvVar("SYNC_COMMITS_INTERVAL_MINUTES", "int").(int),
		SyncUsers:                   getEnvVar("SYNC_USERS_MODE", "bool").(bool),
		SyncUsersIntervalMinutesInt: getEnvVar("SYNC_USERS_INTERVAL_MINUTES", "int").(int),
		GhAccessToken:               getEnvVar("GH_ACCESS_TOKEN", "string").(string),
		TargetLiveGithub:            getEnvVar("TARGET_LIVE_GITHUB", "bool").(bool),
		StartServer:                 getEnvVar("START_SERVER", "bool").(bool),
	}
}

func getEnvVar(name string, expectedType string) interface{} {
	raw := os.Getenv(name)
	if raw == "" {
		logger.LogFatalRedAndExit(name + " is not found in the environment")
	}

	switch expectedType {
	case "string":
		return raw
	case "bool":
		val, err := strconv.ParseBool(raw)
		if err != nil {
			logger.LogFatalRedAndExit(name + " must be a boolean")
		}
		return val
	case "int":
		val, err := strconv.Atoi(raw)
		if err != nil {
			logger.LogFatalRedAndExit(name + " must be an integer")
		}
		return val
	default:
		logger.LogFatalRedAndExit("Unexpected type for " + name)
		return nil
	}
}
