package states

import (
	commands2 "reimagined_eureka/internal/client/cli/commands"
	"reimagined_eureka/internal/client/cli/entities"
)

type CheckMasterKeyState struct {
	GeneralState
}

func NewSignedInState() *CheckMasterKeyState {
	//commands:
	//
	//sign in
	// check login and password locally
	//  if present - check password locally
	//  else - http request
	//sign up
	// http request - login, password
	//  handle errors
	//quit

	cmds := []entities.Command{
		&commands2.QuitCommand{},
	}
	return &CheckMasterKeyState{GeneralState{Commands: cmds}}
}
