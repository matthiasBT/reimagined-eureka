package commands

import (
	"fmt"

	cliEntities "reimagined_eureka/internal/client/cli/entities"
	clientEntities "reimagined_eureka/internal/client/entities"
	"reimagined_eureka/internal/client/infra/logging"
)

type LoginCommand struct {
	Logger          logging.ILogger
	Storage         clientEntities.IStorage
	Proxy           clientEntities.IProxy
	CryptoProvider  clientEntities.ICryptoProvider
	login, password string
}

func (c *LoginCommand) GetName() string {
	return "login"
}

func (c *LoginCommand) GetDescription() string {
	return "log in locally or on server (in case of the first local user's log in operation)"
}

func (c *LoginCommand) Validate(args ...string) error {
	if len(args) != 2 {
		return fmt.Errorf("example: login <login> <password>")
	} // TODO: ******?
	c.login, c.password = args[0], args[1]
	return nil
}

func (c *LoginCommand) Execute() cliEntities.CommandResult {
	user, err := c.Storage.ReadUser(c.login)
	if err != nil {
		return cliEntities.CommandResult{FailureMessage: err.Error()}
	}
	if user != nil {
		err := c.CryptoProvider.VerifyPassword(user, c.password)
		if err != nil {
			msg := fmt.Errorf("password verification failed: %v", err)
			return cliEntities.CommandResult{FailureMessage: msg.Error()}
		}
		return cliEntities.CommandResult{SuccessMessage: "Logged in successfully (locally)"}
	}
	c.Logger.Warningln("User %s not found locally. Going to fetch it from server", c.login)
	_, err = c.Proxy.LogIn(c.login, c.password) // TODO: userData: cookie
	if err != nil {
		msg := fmt.Errorf("failed to log in: %v", err)
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
	return cliEntities.CommandResult{SuccessMessage: "Logged in successfully (on server)"}
}
