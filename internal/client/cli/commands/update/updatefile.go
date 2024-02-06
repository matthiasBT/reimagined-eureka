package update

import (
	"fmt"

	"reimagined_eureka/internal/client/cli/commands"
	"reimagined_eureka/internal/client/cli/commands/add"
	cliEntities "reimagined_eureka/internal/client/cli/entities"
	clientEntities "reimagined_eureka/internal/client/entities"
	"reimagined_eureka/internal/client/infra/logging"
	"reimagined_eureka/internal/common"
)

type UpdateFileCommand struct {
	logger         logging.ILogger
	storage        clientEntities.IStorage
	cryptoProvider clientEntities.ICryptoProvider
	proxy          clientEntities.IProxy
	userID         int
	filePath       string
	rowIDServer    int
	rowIDLocal     int
}

func NewUpdateFileCommand(
	logger logging.ILogger,
	storage clientEntities.IStorage,
	cryptoProvider clientEntities.ICryptoProvider,
	proxy clientEntities.IProxy,
	userID int,
) *UpdateFileCommand {
	return &UpdateFileCommand{
		logger:         logger,
		storage:        storage,
		cryptoProvider: cryptoProvider,
		proxy:          proxy,
		userID:         userID,
	}
}

func (c *UpdateFileCommand) GetName() string {
	return "update-file"
}

func (c *UpdateFileCommand) GetDescription() string {
	return "update binary file contents"
}

func (c *UpdateFileCommand) Validate(args ...string) error {
	if len(args) != 2 {
		return fmt.Errorf("example: update-file <ID> <path>")
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
	c.filePath = args[1]
	return nil
}

func (c *UpdateFileCommand) Execute() cliEntities.CommandResult {
	encrypted, meta, err := add.PrepareFile(c.logger, c.cryptoProvider, c.filePath)
	if err != nil {
		return cliEntities.CommandResult{
			FailureMessage: err.Error(),
		}
	}
	payload := common.FileReq{
		ServerID: &c.rowIDServer,
		Meta:     meta,
		Value:    encrypted,
	}
	if err := c.proxy.UpdateFile(&payload); err != nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("request failed: %v", err).Error(),
		}
	}
	fileLocal := clientEntities.FileLocal{
		File: common.File{
			UserID:           c.userID,
			Meta:             payload.Meta,
			EncryptedContent: payload.Value.Ciphertext,
			Salt:             payload.Value.Salt,
			Nonce:            payload.Value.Nonce,
		},
		ServerID: c.rowIDServer,
	}
	if err := c.storage.SaveFile(&fileLocal); err != nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("failed to update file data locally: %v", err).Error(),
		}
	}
	return cliEntities.CommandResult{
		SuccessMessage: "Successfully updated on server and locally",
	}
}
