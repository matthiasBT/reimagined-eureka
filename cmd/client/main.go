package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

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

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go cli.NewTerminal(logger, storage, serverProxy, cryptoProvider, signals).Run()
	<-signals
	time.Sleep(1 * time.Second)
}
