package states

import (
	commands2 "reimagined_eureka/internal/client/cli/commands"
	"reimagined_eureka/internal/client/cli/entities"
)

type InitialState struct {
	GeneralState
}

func NewInitialState() *InitialState {
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
		&commands2.QuitCommand{},
	}
	return &InitialState{GeneralState{Commands: cmds}}
}
