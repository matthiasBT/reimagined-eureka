package states

import (
	cliCommands "reimagined_eureka/internal/client/cli/commands"
	"reimagined_eureka/internal/client/cli/entities"
	clientEntities "reimagined_eureka/internal/client/entities"
	"reimagined_eureka/internal/client/infra/logging"
)

type InitialState struct {
	GeneralState
	logger  logging.ILogger
	storage clientEntities.IStorage
	proxy   clientEntities.IProxy
}

func NewInitialState(
	logger logging.ILogger,
	storage clientEntities.IStorage,
	proxy clientEntities.IProxy,
	cryptoProvider clientEntities.ICryptoProvider,
) *InitialState {
	//commands:
	//sign in
	// check login and password locally
	//  if present - check password locally
	//  else - http request
	//  check if master_key_checker_encrypted is set (if not - move to SetMasterKeyState)
	//sign up
	// http request - login, password
	//  handle errors
	//  move to SetMasterKeyState
	//quit
	cmds := []entities.Command{
		&cliCommands.LoginCommand{Logger: logger, Storage: storage, Proxy: proxy, CryptoProvider: cryptoProvider},
		// &cliCommands.RegisterCommand{Logger: logger, Storage: storage, Proxy: proxy, CryptoProvider: cryptoProvider},
		&cliCommands.QuitCommand{},
	}
	return &InitialState{GeneralState: GeneralState{Commands: cmds}, logger: logger, storage: storage, proxy: proxy}
}
