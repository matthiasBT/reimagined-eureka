package states

import (
	cliCommands "reimagined_eureka/internal/client/cli/commands"
	"reimagined_eureka/internal/client/cli/commands/add"
	delete2 "reimagined_eureka/internal/client/cli/commands/delete"
	"reimagined_eureka/internal/client/cli/commands/reveal"
	"reimagined_eureka/internal/client/cli/commands/update"
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
	userID int,
	masterKey string,
) *PrimaryState {
	cmds := createCommands(logger, storage, cryptoProvider, proxy, login, password, userID)
	return &PrimaryState{
		GeneralState:   GeneralState{Commands: cmds},
		logger:         logger,
		storage:        storage,
		cryptoProvider: cryptoProvider,
		proxy:          proxy,
		login:          login,
		password:       password,
		userID:         userID,
		masterKey:      masterKey,
	}
}

func (s *PrimaryState) Execute(line string) (entities.State, entities.CommandResult) {
	state, result := s.GeneralState.Execute(line)
	if result.Quit {
		return state, result
	}
	if result.SessionCookie != "" {
		s.proxy.SetSessionCookie(result.SessionCookie)
		s.Commands = createCommands(
			s.logger,
			s.storage,
			s.cryptoProvider,
			s.proxy,
			s.login,
			s.password,
			s.userID,
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
	userID int,
) []entities.Command {
	return []entities.Command{
		cliCommands.NewRefreshSessionCommand(logger, proxy, login, password),

		cliCommands.NewListSecretsCommand(logger, storage, cryptoProvider, userID),

		add.NewAddCredsCommand(logger, storage, cryptoProvider, proxy, userID),
		add.NewAddNoteCommand(logger, storage, cryptoProvider, proxy, userID),
		add.NewAddFileCommand(logger, storage, cryptoProvider, proxy, userID),
		add.NewAddCardCommand(logger, storage, cryptoProvider, proxy, userID),

		reveal.NewRevealCredsCommand(logger, storage, cryptoProvider, userID),
		reveal.NewRevealNoteCommand(logger, storage, cryptoProvider, userID),
		reveal.NewRevealFileCommand(logger, storage, cryptoProvider, userID),
		reveal.NewRevealCardCommand(logger, storage, cryptoProvider, userID),

		update.NewUpdateCredsCommand(logger, storage, cryptoProvider, proxy, userID),
		update.NewUpdateNoteCommand(logger, storage, cryptoProvider, proxy, userID),
		update.NewUpdateFileCommand(logger, storage, cryptoProvider, proxy, userID),
		update.NewUpdateCardCommand(logger, storage, cryptoProvider, proxy, userID),

		delete2.NewDeleteCredsCommand(storage, proxy, userID),
		delete2.NewDeleteNoteCommand(storage, proxy, userID),
		delete2.NewDeleteFileCommand(storage, proxy, userID),

		&cliCommands.QuitCommand{},
	}
}
