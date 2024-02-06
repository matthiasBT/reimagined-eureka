package commands

import (
	"fmt"

	cliEntities "reimagined_eureka/internal/client/cli/entities"
	clientEntities "reimagined_eureka/internal/client/entities"
	"reimagined_eureka/internal/client/infra/logging"
	"reimagined_eureka/internal/common"
)

type UpdateNoteCommand struct {
	Logger         logging.ILogger
	Storage        clientEntities.IStorage
	CryptoProvider clientEntities.ICryptoProvider
	proxy          clientEntities.IProxy
	userID         int
	rowIDServer    int
	rowIDLocal     int
}

func NewUpdateNoteCommand(
	logger logging.ILogger,
	storage clientEntities.IStorage,
	cryptoProvider clientEntities.ICryptoProvider,
	proxy clientEntities.IProxy,
	userID int,
) *UpdateNoteCommand {
	return &UpdateNoteCommand{
		Logger:         logger,
		Storage:        storage,
		CryptoProvider: cryptoProvider,
		proxy:          proxy,
		userID:         userID,
	}
}

func (c *UpdateNoteCommand) GetName() string {
	return "update-note"
}

func (c *UpdateNoteCommand) GetDescription() string {
	return "update a text note"
}

func (c *UpdateNoteCommand) Validate(args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("example: update-note <ID>")
	}
	rowID, err := parsePositiveInt(args[0])
	if err != nil {
		return err
	}
	note, err := c.Storage.ReadNote(c.userID, rowID)
	if err != nil {
		return fmt.Errorf("failed to read note: %v", err)
	}
	if note == nil {
		return fmt.Errorf("note %d doesn't exist", rowID)
	}
	c.rowIDServer = note.ServerID
	c.rowIDLocal = note.ServerID
	return nil
}

func (c *UpdateNoteCommand) Execute() cliEntities.CommandResult {
	encrypted, meta, err := prepareNote(c.Logger, c.CryptoProvider)
	if err != nil {
		return cliEntities.CommandResult{
			FailureMessage: err.Error(),
		}
	}
	payload := common.NoteReq{
		ServerID: &c.rowIDServer,
		Meta:     meta,
		Value:    encrypted,
	}
	if err := c.proxy.UpdateNote(&payload); err != nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("request failed: %v", err).Error(),
		}
	}
	noteLocal := clientEntities.NoteLocal{
		Note: common.Note{
			UserID:           c.userID,
			Meta:             payload.Meta,
			EncryptedContent: payload.Value.Ciphertext,
			Salt:             payload.Value.Salt,
			Nonce:            payload.Value.Nonce,
		},
		ServerID: c.rowIDServer,
	}
	if err := c.Storage.SaveNote(&noteLocal); err != nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("failed to update note locally: %v", err).Error(),
		}
	}
	return cliEntities.CommandResult{
		SuccessMessage: "Successfully updated on server and locally",
	}
}
