package states

import (
	commands2 "reimagined_eureka/internal/client/cli/commands"
	"reimagined_eureka/internal/client/cli/entities"
)

type InitialState struct {
	GeneralState
}

func NewInitialState() *InitialState {
	cmds := []entities.Command{
		&commands2.QuitCommand{},
		&commands2.AddCommand{},
	}
	return &InitialState{GeneralState{Commands: cmds}}
}
