package states

import (
	commands2 "reimagined_eureka/internal/client/cli/commands"
	"reimagined_eureka/internal/client/cli/entities"
)

type SignedInState struct {
	GeneralState
}

func NewSignedInState() *SignedInState {
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
	return &SignedInState{GeneralState{Commands: cmds}}
}
