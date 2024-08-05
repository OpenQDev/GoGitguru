package setup

import (
	"os"
	"strconv"

	"github.com/OpenQDev/GoGitguru/util/logger"
	"github.com/joho/godotenv"
)

type EnvConfig struct {
	PortString                   string
	DbUrl                        string
	OriginUrl                    string
	Debug                        bool
	GhAccessToken                string
	GhAccessTokens               string
	TargetLiveGithub             bool
	RepoSyncInterval             int
	UserSyncInterval             int
	UserDependenciesSyncInterval int
	GitguruUrl                   string
	GitguruApiKey                string
	Environment                  string
	RepoSyncConsumerCount        int
	RepoUrlsConsumerGroup        string
	RepoUrlsTopic                string
	UserSyncConsumerCount        int
	UserSyncConsumerGroup        string
	UserSyncTopic                string
	UserDepsSyncTopic            string
	UserDepsSyncConsumerGroup    string

	KafkaBrokerUrls string
	AwsRegion       string
}

func ExtractAndVerifyEnvironment(pathToDotenv string) EnvConfig {
	if os.Getenv("APP_ENV") != "production" {
		godotenv.Load(pathToDotenv)
	}

	return EnvConfig{
		PortString:                   getEnvVar("PORT", "string").(string),
		DbUrl:                        getEnvVar("DB_URL", "string").(string),
		OriginUrl:                    getEnvVar("ORIGIN_URL", "string").(string),
		Debug:                        getEnvVar("DEBUG_MODE", "bool").(bool),
		GhAccessTokens:               getEnvVar("GH_ACCESS_TOKENS", "string").(string),
		TargetLiveGithub:             getEnvVar("TARGET_LIVE_GITHUB", "bool").(bool),
		RepoSyncInterval:             getEnvVar("REPOSYNC_INTERVAL", "int").(int),
		UserSyncInterval:             getEnvVar("USERSYNC_INTERVAL", "int").(int),
		UserDependenciesSyncInterval: getEnvVar("USERSYNC_DEPDENCIES_INTERVAL", "int").(int),
		GitguruUrl:                   getEnvVar("GITGURU_URL", "string").(string),
		GitguruApiKey:                getEnvVar("API_KEY", "string").(string),
		Environment:                  getEnvVar("APP_ENV", "string").(string),
		RepoSyncConsumerCount:        getEnvVar("REPO_SYNC_CONSUMER_COUNT", "int").(int),
		RepoUrlsConsumerGroup:        getEnvVar("REPO_URLS_CONSUMER_GROUP", "string").(string),
		RepoUrlsTopic:                getEnvVar("REPO_URLS_TOPIC", "string").(string),
		UserSyncConsumerCount:        getEnvVar("USER_SYNC_CONSUMER_COUNT", "int").(int),
		UserSyncConsumerGroup:        getEnvVar("USER_SYNC_CONSUMER_GROUP", "string").(string),
		UserDepsSyncTopic:            getEnvVar("USER_DEPS_SYNC_TOPIC", "string").(string),
		UserDepsSyncConsumerGroup:    getEnvVar("USER_DEPS_SYNC_CONSUMER_GROUP", "string").(string),
		UserSyncTopic:                getEnvVar("USER_SYNC_TOPIC", "string").(string),
		KafkaBrokerUrls:              getEnvVar("KAFKA_BROKER_URLS", "string").(string),
		AwsRegion:                    getEnvVar("AWS_REGION", "string").(string),
	}
}

func getEnvVar(name string, expectedType string) interface{} {
	raw := os.Getenv(name)
	if raw == "" {
		logger.LogFatalRedAndExit("%s is not found in the environment", name)
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
