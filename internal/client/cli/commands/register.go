package commands

import (
	"fmt"

	cliEntities "reimagined_eureka/internal/client/cli/entities"
	clientEntities "reimagined_eureka/internal/client/entities"
)

type RegisterCommand struct {
	Storage         clientEntities.IStorage
	Proxy           clientEntities.IProxy
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
	// TODO: try to query the database
	_, err := c.Proxy.LogIn(c.login, c.password)
	if err != nil {
		msg := fmt.Sprintf("Failed to sign up: %v", err)
		return cliEntities.CommandResult{FailureMessage: msg}
	}
	return cliEntities.CommandResult{SuccessMessage: "Registration successful"}
}
