package delete

import (
	"fmt"

	"reimagined_eureka/internal/client/cli/commands"
	cliEntities "reimagined_eureka/internal/client/cli/entities"
	clientEntities "reimagined_eureka/internal/client/entities"
)

type DeleteFileCommand struct {
	storage     clientEntities.IStorage
	proxy       clientEntities.IProxy
	userID      int
	rowIDServer int
	rowIDLocal  int
}

func NewDeleteFileCommand(
	storage clientEntities.IStorage,
	proxy clientEntities.IProxy,
	userID int,
) *DeleteFileCommand {
	return &DeleteFileCommand{
		storage: storage,
		proxy:   proxy,
		userID:  userID,
	}
}

func (c *DeleteFileCommand) GetName() string {
	return "delete-file"
}

func (c *DeleteFileCommand) GetDescription() string {
	return "delete a file"
}

func (c *DeleteFileCommand) Validate(args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("example: delete-file <ID>")
	}
	rowID, err := commands.ParsePositiveInt(args[0])
	if err != nil {
		return err
	}
	file, err := c.storage.ReadFile(c.userID, rowID)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}
	if file == nil {
		return fmt.Errorf("file %d doesn't exist", rowID)
	}
	c.rowIDServer = file.ServerID
	c.rowIDLocal = file.ID
	return nil
}

func (c *DeleteFileCommand) Execute() cliEntities.CommandResult {
	if err := c.proxy.DeleteFile(c.rowIDServer); err != nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("request failed: %v", err).Error(),
		}
	}
	if err := c.storage.DeleteFile(c.rowIDLocal); err != nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("failed to delete file locally: %v", err).Error(),
		}
	}
	return cliEntities.CommandResult{
		SuccessMessage: "Successfully deleted on server and locally",
	}
}
