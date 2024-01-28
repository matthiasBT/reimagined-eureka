package commands

import (
	"fmt"

	cliEntities "reimagined_eureka/internal/client/cli/entities"
	clientEntities "reimagined_eureka/internal/client/entities"
	"reimagined_eureka/internal/client/infra/logging"
)

type RegisterCommand struct {
	LoginCommand
	Logger          logging.ILogger
	Storage         clientEntities.IStorage
	Proxy           clientEntities.IProxy
	CryptoProvider  clientEntities.ICryptoProvider
	login, password string
}

func NewRegisterCommand(
	logger logging.ILogger,
	storage clientEntities.IStorage,
	proxy clientEntities.IProxy,
	cryptoProvider clientEntities.ICryptoProvider,
) *RegisterCommand {
	return &RegisterCommand{
		LoginCommand: LoginCommand{
			Logger:         logger,
			Storage:        storage,
			Proxy:          proxy,
			CryptoProvider: cryptoProvider,
		},
	}
}

func (c *RegisterCommand) GetName() string {
	return "register"
}

func (c *RegisterCommand) GetDescription() string {
	return "register a new user (requires a server connection)"
}

func (c *RegisterCommand) Validate(args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("example: register <login>")
	}
	password, err := c.readPasswordMasked()
	if err != nil {
		return fmt.Errorf("failed to read password: %v", err)
	}
	c.LoginCommand.login, c.LoginCommand.password = args[0], password
	return nil
}

func (c *RegisterCommand) Execute() cliEntities.CommandResult {
	userData, err := c.LoginCommand.Proxy.Register(c.LoginCommand.login, c.LoginCommand.password)
	if err != nil {
		msg := fmt.Errorf("failed to sign up: %v", err)
		return cliEntities.CommandResult{FailureMessage: msg.Error()}
	}
	newUser := &clientEntities.User{Login: c.LoginCommand.login}
	if err := c.LoginCommand.CryptoProvider.HashPassword(newUser, c.LoginCommand.password); err != nil {
		msg := fmt.Errorf("failed to store user %s data locally: %v", newUser.Login, err)
		return cliEntities.CommandResult{FailureMessage: msg.Error()}
	}
	if err := c.LoginCommand.Storage.SaveUser(newUser); err != nil {
		msg := fmt.Errorf("failed to store user %s data locally: %v", newUser.Login, err)
		return cliEntities.CommandResult{FailureMessage: msg.Error()}
	}
	return cliEntities.CommandResult{
		SuccessMessage: "Registered successfully",
		SessionCookie:  userData.SessionCookie,
	}
}
