package main

import (
	"reimagined_eureka/internal/client/adapters"
	"reimagined_eureka/internal/client/infra/config"
	"reimagined_eureka/internal/client/infra/logging"
)

func main() {
	logger := logging.SetupLogger()
	conf, err := config.InitConfig()
	if err != nil {
		logger.Error("Failed to start client: %v", err)
	}
	logger.Success("Client config loaded")
	_, err = adapters.NewSQLiteStorage(logger, conf.DatabasePath)
	if err != nil {
		logger.Error("Failed to start client: %v", err)
	}
	logger.Success("Database initialized")
	//t := cli.NewTerminal()
	//t.Run()
}
