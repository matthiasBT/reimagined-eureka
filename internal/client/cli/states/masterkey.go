package states

import (
	cliCommands "reimagined_eureka/internal/client/cli/commands"
	cliEntities "reimagined_eureka/internal/client/cli/entities"
	clientEntities "reimagined_eureka/internal/client/entities"
	"reimagined_eureka/internal/client/infra/logging"
)

type MasterKeyState struct {
	GeneralState
	logger         logging.ILogger
	storage        clientEntities.IStorage
	cryptoProvider clientEntities.ICryptoProvider
	login          string
	userID         int
}

func NewMasterKeyState(
	logger logging.ILogger,
	storage clientEntities.IStorage,
	cryptoProvider clientEntities.ICryptoProvider,
	login string,
	userID int,
) *MasterKeyState {
	cmds := []cliEntities.Command{
		cliCommands.NewMasterKeyCommand(logger, storage, cryptoProvider, login),
		&cliCommands.QuitCommand{},
	}
	return &MasterKeyState{
		GeneralState:   GeneralState{Commands: cmds},
		logger:         logger,
		storage:        storage,
		cryptoProvider: cryptoProvider,
		login:          login,
		userID:         userID,
	}
}

func (s *MasterKeyState) Execute(line string) (cliEntities.State, cliEntities.CommandResult) {
	state, result := s.GeneralState.Execute(line)
	if result.MasterKey != "" {
		state = NewPrimaryState(s.logger, s.storage, s.cryptoProvider, s.userID, result.MasterKey)
	}
	return state, result
}
