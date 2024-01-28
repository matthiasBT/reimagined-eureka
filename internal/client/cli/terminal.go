package cli

import (
	"bufio"
	"os"

	"reimagined_eureka/internal/client/adapters"
	cliEntities "reimagined_eureka/internal/client/cli/entities"
	cliStates "reimagined_eureka/internal/client/cli/states"
	clientEntities "reimagined_eureka/internal/client/entities"
	"reimagined_eureka/internal/client/infra/logging"
)

const PromptGeneric = "> "
const PromptOnExit = "Bye!"

type Terminal struct {
	logger    logging.ILogger
	storage   clientEntities.IStorage
	proxy     clientEntities.IProxy
	currState cliEntities.State
	scanner   *bufio.Scanner
}

func NewTerminal(
	logger logging.ILogger,
	storage clientEntities.IStorage,
	proxy clientEntities.IProxy,
	cryptoProvider *adapters.CryptoProvider,
) *Terminal {
	scanner := bufio.NewScanner(os.Stdin) // TODO: pass as parameter
	return &Terminal{
		logger:    logger,
		storage:   storage,
		proxy:     proxy,
		currState: cliStates.NewInitialState(logger, storage, proxy, cryptoProvider),
		scanner:   scanner,
	}
}

func (t *Terminal) Run() {
	for {
		t.logger.Infoln(t.currState.GetPrompt())
		t.logger.Info(PromptGeneric)
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			input := scanner.Text()
			nextState, result := t.currState.Execute(input)
			if result.SuccessMessage != "" {
				t.logger.Successln(result.SuccessMessage)
			} else if result.FailureMessage != "" {
				t.logger.Failureln(result.FailureMessage)
			}
			if nextState != nil {
				if _, ok := nextState.(*cliStates.QuitState); ok {
					t.logger.Successln(PromptOnExit)
					return
				}
				t.currState = nextState
			}
			t.logger.Infoln("")
			continue
		}
		t.handleInputErrors()
		t.logger.Successln(PromptOnExit)
		return
	}
}

func (t *Terminal) handleInputErrors() {
	if err := t.scanner.Err(); err != nil {
		t.logger.Failureln("Failureln reading input:", err) // TODO: handle errors
		os.Exit(1)
	}
}
