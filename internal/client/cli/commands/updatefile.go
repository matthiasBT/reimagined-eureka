package commands

import (
	"fmt"

	cliEntities "reimagined_eureka/internal/client/cli/entities"
	clientEntities "reimagined_eureka/internal/client/entities"
	"reimagined_eureka/internal/client/infra/logging"
	"reimagined_eureka/internal/common"
)

type UpdateFileCommand struct {
	Logger         logging.ILogger
	Storage        clientEntities.IStorage
	CryptoProvider clientEntities.ICryptoProvider
	proxy          clientEntities.IProxy
	userID         int
	filePath       string
	rowID          int
}

func NewUpdateFileCommand(
	logger logging.ILogger,
	storage clientEntities.IStorage,
	cryptoProvider clientEntities.ICryptoProvider,
	proxy clientEntities.IProxy,
	userID int,
) *UpdateFileCommand {
	return &UpdateFileCommand{
		Logger:         logger,
		Storage:        storage,
		CryptoProvider: cryptoProvider,
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
	rowID, err := parsePositiveInt(args[0])
	if err != nil {
		return err
	}
	note, err := c.Storage.ReadFile(c.userID, rowID)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}
	if note == nil {
		return fmt.Errorf("file %d doesn't exist", rowID)
	}
	c.rowID = rowID
	c.filePath = args[1]
	return nil
}

func (c *UpdateFileCommand) Execute() cliEntities.CommandResult {
	encrypted, meta, err := prepareFile(c.Logger, c.CryptoProvider, c.filePath)
	if err != nil {
		return cliEntities.CommandResult{
			FailureMessage: err.Error(),
		}
	}
	payload := common.FileReq{
		ServerID: &c.rowID,
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
		ServerID: c.rowID,
	}
	if err := c.Storage.SaveFile(&fileLocal); err != nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("failed to update file data locally: %v", err).Error(),
		}
	}
	return cliEntities.CommandResult{
		SuccessMessage: "Successfully updated on server and locally",
	}
}
