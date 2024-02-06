package states

import (
	"fmt"
	"strings"

	cliEntities "reimagined_eureka/internal/client/cli/entities"
)

const ResultUnknownCommand = "Unknown command"

type GeneralState struct {
	Commands []cliEntities.Command
}

func (s GeneralState) GetPrompt() string {
	var result = []string{"Please type one of the following commands:"}
	for _, cmd := range s.Commands {
		cmdPrompt := fmt.Sprintf("• %s: %s", cmd.GetName(), cmd.GetDescription())
		result = append(result, cmdPrompt)
	}
	return strings.Join(result, "\n")
}

func (s GeneralState) Execute(line string) (cliEntities.State, cliEntities.CommandResult) {
	parts := strings.Fields(line)
	for _, cmd := range s.Commands {
		if parts[0] == cmd.GetName() {
			if err := cmd.Validate(parts[1:]...); err != nil {
				return nil, cliEntities.CommandResult{FailureMessage: err.Error()}
			}
			result := cmd.Execute()
			if result.Quit {
				return &QuitState{}, result
			}
			return nil, result
		}
	}
	return nil, cliEntities.CommandResult{FailureMessage: ResultUnknownCommand}
}
