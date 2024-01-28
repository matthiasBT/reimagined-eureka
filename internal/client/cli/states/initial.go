package states

import (
	cliCommands "reimagined_eureka/internal/client/cli/commands"
	"reimagined_eureka/internal/client/cli/entities"
	clientEntities "reimagined_eureka/internal/client/entities"
)

type InitialState struct {
	GeneralState
	storage clientEntities.IStorage
	proxy   clientEntities.IProxy
}

func NewInitialState(storage clientEntities.IStorage, proxy clientEntities.IProxy) *InitialState {
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
		&cliCommands.LoginCommand{Storage: storage, Proxy: proxy},
		&cliCommands.QuitCommand{},
	}
	return &InitialState{GeneralState{Commands: cmds}, storage, proxy}
}
