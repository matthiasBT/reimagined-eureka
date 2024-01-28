package commands

import (
	"fmt"

	cliEntities "reimagined_eureka/internal/client/cli/entities"
	clientEntities "reimagined_eureka/internal/client/entities"
)

type LoginCommand struct {
	Storage         clientEntities.IStorage
	Proxy           clientEntities.IProxy
	login, password string
}

func (c *LoginCommand) GetName() string {
	return "login"
}

func (c *LoginCommand) GetDescription() string {
	return "log in locally or on server (in case of the first local user login)"
}

func (c *LoginCommand) Validate(args ...string) error {
	if len(args) != 2 {
		return fmt.Errorf("example: login <login> <password>")
	} // TODO: ******?
	c.login, c.password = args[0], args[1]
	return nil
}

func (c *LoginCommand) Execute() cliEntities.CommandResult {
	// TODO: try to query the database
	_, err := c.Proxy.LogIn(c.login, c.password)
	if err != nil {
		msg := fmt.Sprintf("Failed to log in: %v", err)
		return cliEntities.CommandResult{FailureMessage: msg}
	}
	return cliEntities.CommandResult{SuccessMessage: "Logged in successfully"}
}
