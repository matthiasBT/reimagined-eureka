package states

import (
	cliCommands "reimagined_eureka/internal/client/cli/commands"
	"reimagined_eureka/internal/client/cli/entities"
	clientEntities "reimagined_eureka/internal/client/entities"
	"reimagined_eureka/internal/client/infra/logging"
)

type PrimaryState struct {
	GeneralState
	logger         logging.ILogger
	storage        clientEntities.IStorage
	cryptoProvider clientEntities.ICryptoProvider
	login          string
	masterKey      string
}

func NewPrimaryState(
	logger logging.ILogger,
	storage clientEntities.IStorage,
	cryptoProvider clientEntities.ICryptoProvider,
	login string,
	masterKey string,
) *PrimaryState {
	cmds := []entities.Command{
		cliCommands.NewListSecretsCommand(logger, storage, cryptoProvider, login),
		&cliCommands.QuitCommand{},
	}
	return &PrimaryState{
		GeneralState:   GeneralState{Commands: cmds},
		logger:         logger,
		storage:        storage,
		cryptoProvider: cryptoProvider,
		login:          login,
		masterKey:      masterKey,
	}
}
