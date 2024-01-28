package commands

import (
	"fmt"

	cliEntities "reimagined_eureka/internal/client/cli/entities"
	clientEntities "reimagined_eureka/internal/client/entities"
	"reimagined_eureka/internal/client/infra/logging"
)

type RegisterCommand struct {
	Logger          logging.ILogger
	Storage         clientEntities.IStorage
	Proxy           clientEntities.IProxy
	CryptoProvider  clientEntities.ICryptoProvider
	login, password string
}

func (c *RegisterCommand) GetName() string {
	return "register"
}

func (c *RegisterCommand) GetDescription() string {
	return "register a new user (requires a server connection)"
}

func (c *RegisterCommand) Validate(args ...string) error {
	if len(args) != 2 {
		return fmt.Errorf("example: register <login> <password>") // TODO: ******?
	}
	c.login, c.password = args[0], args[1]
	return nil
}

func (c *RegisterCommand) Execute() cliEntities.CommandResult {
	_, err := c.Proxy.Register(c.login, c.password) // TODO: userData: cookie
	if err != nil {
		msg := fmt.Errorf("failed to sign up: %v", err)
		return cliEntities.CommandResult{FailureMessage: msg.Error()}
	}
	newUser := &clientEntities.User{Login: c.login}
	if err := c.CryptoProvider.HashPassword(newUser, c.password); err != nil {
		msg := fmt.Errorf("failed to store user %s data locally: %v", newUser.Login, err)
		return cliEntities.CommandResult{FailureMessage: msg.Error()}
	}
	if err := c.Storage.SaveUser(newUser); err != nil {
		msg := fmt.Errorf("failed to store user %s data locally: %v", newUser.Login, err)
		return cliEntities.CommandResult{FailureMessage: msg.Error()}
	}
	return cliEntities.CommandResult{SuccessMessage: "Registered successfully"}
}
