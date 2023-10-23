package main

import (
	"server"
	"util/logger"
	"util/setup"
)

func main() {
	env := setup.ExtractAndVerifyEnvironment(".env")

	database, apiCfg := server.PrepareServerSingleton(env.DbUrl)
	logger.SetDebugMode(env.Debug)
	server.StartServer(apiCfg, env.PortString, env.OriginUrl)
}
