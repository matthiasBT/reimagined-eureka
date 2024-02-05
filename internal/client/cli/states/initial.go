package states

import (
	cliCommands "reimagined_eureka/internal/client/cli/commands"
	cliEntities "reimagined_eureka/internal/client/cli/entities"
	clientEntities "reimagined_eureka/internal/client/entities"
	"reimagined_eureka/internal/client/infra/logging"
)

type InitialState struct {
	GeneralState
	logger         logging.ILogger
	storage        clientEntities.IStorage
	proxy          clientEntities.IProxy
	cryptoProvider clientEntities.ICryptoProvider
}

func NewInitialState(
	logger logging.ILogger,
	storage clientEntities.IStorage,
	proxy clientEntities.IProxy,
	cryptoProvider clientEntities.ICryptoProvider,
) *InitialState {
	cmds := []cliEntities.Command{
		&cliCommands.LoginCommand{
			Logger:         logger,
			Storage:        storage,
			Proxy:          proxy,
			CryptoProvider: cryptoProvider,
		},
		cliCommands.NewRegisterCommand(logger, storage, proxy, cryptoProvider),
		&cliCommands.QuitCommand{},
	}
	return &InitialState{
		GeneralState:   GeneralState{Commands: cmds},
		logger:         logger,
		storage:        storage,
		proxy:          proxy,
		cryptoProvider: cryptoProvider,
	}
}

func (s InitialState) Execute(line string) (cliEntities.State, cliEntities.CommandResult) {
	state, result := s.GeneralState.Execute(line)
	if result.Login != "" {
		state = NewMasterKeyState(
			s.logger,
			s.storage,
			s.cryptoProvider,
			s.proxy,
			result.Login,
			result.Password,
			result.SessionCookie,
			result.UserID,
		)
	}
	return state, result
}
