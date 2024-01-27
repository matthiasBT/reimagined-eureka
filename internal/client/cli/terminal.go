package cli

import (
	"bufio"
	"fmt"
	"os"

	"reimagined_eureka/internal/client/cli/entities"
	states2 "reimagined_eureka/internal/client/cli/states"
)

const PromptGeneric = "> "
const PromptOnExit = "Bye!"

type Terminal struct {
	currState entities.State
	scanner   *bufio.Scanner
}

func NewTerminal() *Terminal {
	return &Terminal{
		currState: states2.NewInitialState(),
		scanner:   bufio.NewScanner(os.Stdin), // TODO: pass as parameter
	}
}

func (t *Terminal) Run() {
	for {
		fmt.Println(t.currState.GetPrompt())
		fmt.Print(PromptGeneric)
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			input := scanner.Text()
			nextState, result := t.currState.Execute(input)
			if result.SuccessMessage != "" {
				fmt.Println(result.SuccessMessage)
			} else if result.FailureMessage != "" {
				fmt.Println(result.FailureMessage)
			}
			if nextState != nil {
				if _, ok := nextState.(*states2.QuitState); ok {
					t.exitNormal()
				}
				t.currState = nextState
			}
			fmt.Println()
			continue
		}
		t.handleInputErrors()
		t.exitNormal()
	}
}

func (t *Terminal) handleInputErrors() {
	if err := t.scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading input:", err) // TODO: handle errors
		os.Exit(1)
	}
}

func (t *Terminal) exitNormal() {
	fmt.Println(PromptOnExit)
	os.Exit(0)
}
