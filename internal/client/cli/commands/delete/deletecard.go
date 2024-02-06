package delete

import (
	"fmt"

	"reimagined_eureka/internal/client/cli/commands"
	cliEntities "reimagined_eureka/internal/client/cli/entities"
	clientEntities "reimagined_eureka/internal/client/entities"
)

type DeleteCardCommand struct {
	storage     clientEntities.IStorage
	proxy       clientEntities.IProxy
	userID      int
	rowIDServer int
	rowIDLocal  int
}

func NewDeleteCardCommand(
	storage clientEntities.IStorage,
	proxy clientEntities.IProxy,
	userID int,
) *DeleteCardCommand {
	return &DeleteCardCommand{
		storage: storage,
		proxy:   proxy,
		userID:  userID,
	}
}

func (c *DeleteCardCommand) GetName() string {
	return "delete-card"
}

func (c *DeleteCardCommand) GetDescription() string {
	return "delete a card"
}

func (c *DeleteCardCommand) Validate(args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("example: delete-card <ID>")
	}
	rowID, err := commands.ParsePositiveInt(args[0])
	if err != nil {
		return err
	}
	card, err := c.storage.ReadCard(c.userID, rowID)
	if err != nil {
		return fmt.Errorf("failed to read card: %v", err)
	}
	if card == nil {
		return fmt.Errorf("card %d doesn't exist", rowID)
	}
	c.rowIDServer = card.ServerID
	c.rowIDLocal = card.ID
	return nil
}

func (c *DeleteCardCommand) Execute() cliEntities.CommandResult {
	if err := c.proxy.DeleteCard(c.rowIDServer); err != nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("request failed: %v", err).Error(),
		}
	}
	if err := c.storage.DeleteCard(c.rowIDLocal); err != nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("failed to delete card locally: %v", err).Error(),
		}
	}
	return cliEntities.CommandResult{
		SuccessMessage: "Successfully deleted on server and locally",
	}
}
