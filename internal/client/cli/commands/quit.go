package commands

import (
	"awesomeProject1/internal/client/cli/entities"
)

type QuitCommand struct {
}

func (c *QuitCommand) GetName() string {
	return "quit"
}

func (c *QuitCommand) GetDescription() string {
	return "exit the program"
}

func (c *QuitCommand) Validate(args ...string) error {
	return nil // always valid
}

func (c *QuitCommand) Execute() entities.CommandResult {
	return entities.CommandResult{Quit: true}
}
