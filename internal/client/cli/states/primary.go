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
	proxy          clientEntities.IProxy
	login          string
	password       string
	sessionCookie  string
	userID         int
	masterKey      string
}

func NewPrimaryState(
	logger logging.ILogger,
	storage clientEntities.IStorage,
	cryptoProvider clientEntities.ICryptoProvider,
	proxy clientEntities.IProxy,
	login string,
	password string,
	sessionCookie string,
	userID int,
	masterKey string,
) *PrimaryState {
	cmds := createCommands(logger, storage, cryptoProvider, proxy, login, password, sessionCookie, userID, masterKey)
	return &PrimaryState{
		GeneralState:   GeneralState{Commands: cmds},
		logger:         logger,
		storage:        storage,
		cryptoProvider: cryptoProvider,
		proxy:          proxy,
		login:          login,
		password:       password,
		sessionCookie:  sessionCookie,
		userID:         userID,
		masterKey:      masterKey,
	}
}

func (s *PrimaryState) Execute(line string) (entities.State, entities.CommandResult) {
	state, result := s.GeneralState.Execute(line)
	if result.Quit {
		return state, result
	}
	if result.SessionCookie != "" { // TODO: test cookie refresh!
		s.proxy.SetSessionCookie(result.SessionCookie)
		s.Commands = createCommands(
			s.logger,
			s.storage,
			s.cryptoProvider,
			s.proxy,
			s.login,
			s.password,
			result.SessionCookie,
			s.userID,
			s.masterKey,
		)
	}
	return state, result
}

func createCommands(
	logger logging.ILogger,
	storage clientEntities.IStorage,
	cryptoProvider clientEntities.ICryptoProvider,
	proxy clientEntities.IProxy,
	login string,
	password string,
	sessionCookie string,
	userID int,
	masterKey string,
) []entities.Command {
	return []entities.Command{
		cliCommands.NewRefreshSessionCommand(logger, proxy, login, password),
		cliCommands.NewListSecretsCommand(logger, storage, cryptoProvider, userID),
		cliCommands.NewAddCredsCommand(logger, storage, cryptoProvider, proxy, sessionCookie, masterKey, userID),
		cliCommands.NewAddNoteCommand(logger, storage, cryptoProvider, proxy, sessionCookie, masterKey, userID),
		&cliCommands.QuitCommand{},
	}
}
