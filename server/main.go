package main

import (
	server "server/src"
	"util/logger"
	"util/setup"
)

func main() {
	env := setup.ExtractAndVerifyEnvironment(".env")

	_, apiCfg := server.PrepareServerSingleton(env.DbUrl)
	logger.SetDebugMode(env.Debug)
	server.StartServer(apiCfg, env.PortString, env.OriginUrl)
}
