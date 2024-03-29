package main

import (
	server "github.com/OpenQDev/GoGitguru/server/src"
	"github.com/OpenQDev/GoGitguru/util/logger"
	"github.com/OpenQDev/GoGitguru/util/setup"
)

func main() {
	env := setup.ExtractAndVerifyEnvironment("../.env")

	conn, apiCfg := server.PrepareServerSingleton(env.DbUrl)
	defer conn.Close()

	logger.SetDebugMode(env.Debug)
	server.StartServer(apiCfg, env.PortString, env.OriginUrl)
}
