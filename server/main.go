package main

import (
	server "server/src"

	"github.com/OpenQDev/GoGitguru/util/logger"
	"github.com/OpenQDev/GoGitguru/util/setup"
)

func main() {
	env := setup.ExtractAndVerifyEnvironment(".env")

	_, apiCfg := server.PrepareServerSingleton(env.DbUrl)
	logger.SetDebugMode(env.Debug)
	server.StartServer(apiCfg, env.PortString, env.OriginUrl)
}
