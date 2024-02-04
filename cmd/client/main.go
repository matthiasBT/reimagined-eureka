package main

import (
	"reimagined_eureka/internal/client/adapters"
	"reimagined_eureka/internal/client/cli"
	"reimagined_eureka/internal/client/infra/config"
	"reimagined_eureka/internal/client/infra/logging"
)

func main() {
	logger := logging.SetupLogger()

	conf, err := config.InitConfig()
	if err != nil {
		logger.Failureln("Failed to start client: %v", err)
		return
	}
	logger.Successln("Client config loaded")

	storage, err := adapters.NewSQLiteStorage(logger, conf.DatabasePath)
	if err != nil {
		logger.Failureln("Failed to start client: %v", err)
		return
	}
	logger.Successln("Database initialized")
	defer storage.Shutdown()

	serverProxy := adapters.NewServerProxy(conf.ServerURL)
	cryptoProvider := adapters.NewCryptoProvider()

	// TODO: recover from panic to prevent ugly output
	cli.NewTerminal(logger, storage, serverProxy, cryptoProvider).Run()
}
