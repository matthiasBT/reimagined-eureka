package states

import (
	cliCommands "reimagined_eureka/internal/client/cli/commands"
	"reimagined_eureka/internal/client/cli/entities"
	clientEntities "reimagined_eureka/internal/client/entities"
	"reimagined_eureka/internal/client/infra/logging"
)

type MasterKeyState struct {
	GeneralState
	logger         logging.ILogger
	storage        clientEntities.IStorage
	cryptoProvider clientEntities.ICryptoProvider
	login          string
}

func NewMasterKeyState(
	logger logging.ILogger,
	storage clientEntities.IStorage,
	cryptoProvider clientEntities.ICryptoProvider,
	login string,
) *MasterKeyState {
	cmds := []entities.Command{
		cliCommands.NewMasterKeyCommand(logger, storage, cryptoProvider, login),
		&cliCommands.QuitCommand{},
	}
	return &MasterKeyState{
		GeneralState:   GeneralState{Commands: cmds},
		logger:         logger,
		storage:        storage,
		cryptoProvider: cryptoProvider,
		login:          login,
	}
}
