package commands

import (
	"fmt"

	cliEntities "reimagined_eureka/internal/client/cli/entities"
	clientEntities "reimagined_eureka/internal/client/entities"
)

type DeleteNoteCommand struct {
	storage     clientEntities.IStorage
	proxy       clientEntities.IProxy
	userID      int
	rowIDServer int
	rowIDLocal  int
}

func NewDeleteNoteCommand(
	storage clientEntities.IStorage,
	proxy clientEntities.IProxy,
	userID int,
) *DeleteNoteCommand {
	return &DeleteNoteCommand{
		storage: storage,
		proxy:   proxy,
		userID:  userID,
	}
}

func (c *DeleteNoteCommand) GetName() string {
	return "delete-note"
}

func (c *DeleteNoteCommand) GetDescription() string {
	return "delete a note"
}

func (c *DeleteNoteCommand) Validate(args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("example: delete-note <ID>")
	}
	rowID, err := parsePositiveInt(args[0])
	if err != nil {
		return err
	}
	note, err := c.storage.ReadNote(c.userID, rowID)
	if err != nil {
		return fmt.Errorf("failed to read note: %v", err)
	}
	if note == nil {
		return fmt.Errorf("note %d doesn't exist", rowID)
	}
	c.rowIDServer = note.ServerID
	c.rowIDLocal = note.ID
	return nil
}

func (c *DeleteNoteCommand) Execute() cliEntities.CommandResult {
	if err := c.proxy.DeleteNote(c.rowIDServer); err != nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("request failed: %v", err).Error(),
		}
	}
	if err := c.storage.DeleteNote(c.rowIDLocal); err != nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("failed to delete note locally: %v", err).Error(),
		}
	}
	return cliEntities.CommandResult{
		SuccessMessage: "Successfully deleted on server and locally",
	}
}
