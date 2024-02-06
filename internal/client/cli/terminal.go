package cli

import (
	"bufio"
	"os"
	"syscall"

	"reimagined_eureka/internal/client/adapters"
	cliEntities "reimagined_eureka/internal/client/cli/entities"
	cliStates "reimagined_eureka/internal/client/cli/states"
	clientEntities "reimagined_eureka/internal/client/entities"
	"reimagined_eureka/internal/client/infra/logging"
)

const PromptGeneric = "> "
const PromptOnExit = "Bye!"

type Terminal struct {
	logger      logging.ILogger
	storage     clientEntities.IStorage
	proxy       clientEntities.IProxy
	currState   cliEntities.State
	scanner     *bufio.Scanner
	exitChannel chan os.Signal
}

func NewTerminal(
	logger logging.ILogger,
	storage clientEntities.IStorage,
	proxy clientEntities.IProxy,
	cryptoProvider *adapters.CryptoProvider,
	exitChannel chan os.Signal,
) *Terminal {
	scanner := bufio.NewScanner(os.Stdin)
	return &Terminal{
		logger:      logger,
		storage:     storage,
		proxy:       proxy,
		currState:   cliStates.NewInitialState(logger, storage, proxy, cryptoProvider),
		scanner:     scanner,
		exitChannel: exitChannel,
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
					t.exitChannel <- syscall.SIGINT
					return
				}
				t.currState = nextState
			}
			t.logger.Infoln("")
			continue
		}
		t.handleInputErrors()
		t.logger.Successln(PromptOnExit)
		t.exitChannel <- syscall.SIGINT
		return
	}
}

func (t *Terminal) handleInputErrors() {
	if err := t.scanner.Err(); err != nil {
		t.logger.Failureln("Failureln reading input:", err)
		t.exitChannel <- syscall.SIGINT
	}
}
