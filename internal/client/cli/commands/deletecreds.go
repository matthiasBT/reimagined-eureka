package commands

import (
	"fmt"

	cliEntities "reimagined_eureka/internal/client/cli/entities"
	clientEntities "reimagined_eureka/internal/client/entities"
	"reimagined_eureka/internal/client/infra/logging"
)

type DeleteCredsCommand struct {
	Logger      logging.ILogger
	Storage     clientEntities.IStorage
	proxy       clientEntities.IProxy
	userID      int
	rowIDServer int
	rowIDLocal  int
}

func NewDeleteCredsCommand(
	logger logging.ILogger,
	storage clientEntities.IStorage,
	proxy clientEntities.IProxy,
	userID int,
) *DeleteCredsCommand {
	return &DeleteCredsCommand{
		Logger:  logger,
		Storage: storage,
		proxy:   proxy,
		userID:  userID,
	}
}

func (c *DeleteCredsCommand) GetName() string {
	return "delete-creds"
}

func (c *DeleteCredsCommand) GetDescription() string {
	return "delete a login-password pair"
}

func (c *DeleteCredsCommand) Validate(args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("example: delete-creds <ID>")
	}
	rowID, err := parsePositiveInt(args[0])
	if err != nil {
		return err
	}
	creds, err := c.Storage.ReadCredential(c.userID, rowID)
	if err != nil {
		return fmt.Errorf("failed to read creds: %v", err)
	}
	if creds == nil {
		return fmt.Errorf("creds %d don't exist", rowID)
	}
	c.rowIDServer = creds.ServerID
	c.rowIDLocal = creds.ID
	return nil
}

func (c *DeleteCredsCommand) Execute() cliEntities.CommandResult {
	if err := c.proxy.DeleteCredentials(c.rowIDServer); err != nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("request failed: %v", err).Error(),
		}
	}
	if err := c.Storage.DeleteCredentials(c.rowIDLocal); err != nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("failed to delete credentials locally: %v", err).Error(),
		}
	}
	return cliEntities.CommandResult{
		SuccessMessage: "Successfully deleted on server and locally",
	}
}
