package states

import (
	entities2 "awesomeProject1/internal/client/cli/entities"
	"fmt"
	"strings"
)

const ResultUnknownCommand = "Unknown command"

type GeneralState struct {
	Commands []entities2.Command
}

func (s GeneralState) GetPrompt() string {
	var result = []string{"Please type one of the following commands:"}
	for _, cmd := range s.Commands {
		cmdPrompt := fmt.Sprintf("• %s: %s", cmd.GetName(), cmd.GetDescription())
		result = append(result, cmdPrompt)
	}
	return strings.Join(result, "\n")
}

func (s GeneralState) Execute(line string) (entities2.State, entities2.CommandResult) {
	parts := strings.Fields(line)
	for _, cmd := range s.Commands {
		if parts[0] == cmd.GetName() {
			if err := cmd.Validate(parts[1:]...); err != nil {
				return nil, entities2.CommandResult{FailureMessage: err.Error()}
			}
			result := cmd.Execute()
			if result.Quit {
				return &QuitState{}, result
			}
			return nil, result // no new state as a result of command
		}
	}
	return nil, entities2.CommandResult{FailureMessage: ResultUnknownCommand}
}
