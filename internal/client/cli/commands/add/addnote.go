package add

import (
	"fmt"

	"reimagined_eureka/internal/client/cli/commands"
	cliEntities "reimagined_eureka/internal/client/cli/entities"
	clientEntities "reimagined_eureka/internal/client/entities"
	"reimagined_eureka/internal/client/infra/logging"
	"reimagined_eureka/internal/common"
)

type AddNoteCommand struct {
	logger         logging.ILogger
	storage        clientEntities.IStorage
	cryptoProvider clientEntities.ICryptoProvider
	proxy          clientEntities.IProxy
	userID         int
}

func NewAddNoteCommand(
	logger logging.ILogger,
	storage clientEntities.IStorage,
	cryptoProvider clientEntities.ICryptoProvider,
	proxy clientEntities.IProxy,
	userID int,
) *AddNoteCommand {
	return &AddNoteCommand{
		logger:         logger,
		storage:        storage,
		cryptoProvider: cryptoProvider,
		proxy:          proxy,
		userID:         userID,
	}
}

func (c *AddNoteCommand) GetName() string {
	return "add-note"
}

func (c *AddNoteCommand) GetDescription() string {
	return "add a text note"
}

func (c *AddNoteCommand) Validate(args ...string) error {
	if len(args) != 0 {
		return fmt.Errorf("example: add-note")
	}
	return nil
}

func (c *AddNoteCommand) Execute() cliEntities.CommandResult {
	encrypted, meta, err := PrepareNote(c.logger, c.cryptoProvider)
	if err != nil {
		return cliEntities.CommandResult{
			FailureMessage: err.Error(),
		}
	}
	payload := common.NoteReq{
		ServerID: nil,
		Meta:     meta,
		Value:    encrypted,
	}
	rowID, err := c.proxy.AddNote(&payload)
	if err != nil {
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
		ServerID: rowID,
	}
	if err := c.storage.SaveNote(&noteLocal); err != nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("failed to store note locally: %v", err).Error(),
		}
	}
	return cliEntities.CommandResult{
		SuccessMessage: "Successfully stored on server and locally",
	}
}

func PrepareNote(
	logger logging.ILogger, cryptoProvider clientEntities.ICryptoProvider,
) (*common.EncryptionResult, string, error) {
	content, err := commands.ReadNonSecretValue(logger, "note content") // don't replace with '*' because it's multiline
	if err != nil {
		return nil, "", fmt.Errorf("failed to read note content: %v", err)
	}
	meta, err := commands.ReadNonSecretValue(logger, "meta information")
	if err != nil {
		return nil, "", fmt.Errorf("failed to read meta information: %v", err)
	}
	encrypted, err := cryptoProvider.Encrypt([]byte(content))
	if err != nil {
		return nil, "", fmt.Errorf("failed to encrypt note content: %v", err)
	}
	return encrypted, meta, nil
}
