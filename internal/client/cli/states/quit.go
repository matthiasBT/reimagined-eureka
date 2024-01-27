package states

import (
	entities2 "awesomeProject1/internal/client/cli/entities"
)

type QuitState struct {
	GeneralState
}

const ErrMessage = "QuitState is not a regular state, it shouldn't be used"

func (s *QuitState) GetPrompt() string {
	panic(ErrMessage)
}

func (s *QuitState) Execute(string) (entities2.State, entities2.CommandResult) {
	panic(ErrMessage)
}
