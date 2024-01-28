package main

import (
	"reimagined_eureka/internal/client/adapters"
	"reimagined_eureka/internal/client/cli"
	"reimagined_eureka/internal/client/infra/config"
	"reimagined_eureka/internal/client/infra/logging"
)

func main() {
	logger := logging.SetupLogger()

	// config init
	conf, err := config.InitConfig()
	if err != nil {
		logger.Failureln("Failed to start client: %v", err)
		return
	}
	logger.Successln("Client config loaded")

	// storage init
	storage, err := adapters.NewSQLiteStorage(logger, conf.DatabasePath)
	if err != nil {
		logger.Failureln("Failed to start client: %v", err)
		return
	}
	logger.Successln("Database initialized")
	defer storage.Shutdown()

	// proxy init
	serverProxy := adapters.NewServerProxy(conf.ServerURL)

	// launch
	cli.NewTerminal(logger, storage, serverProxy).Run()
}
