package commands

import (
	"fmt"
	"io"
	"os"

	cliEntities "reimagined_eureka/internal/client/cli/entities"
	clientEntities "reimagined_eureka/internal/client/entities"
	"reimagined_eureka/internal/client/infra/logging"
	"reimagined_eureka/internal/common"
)

type AddFileCommand struct {
	Logger         logging.ILogger
	Storage        clientEntities.IStorage
	CryptoProvider clientEntities.ICryptoProvider
	proxy          clientEntities.IProxy
	userID         int
	filePath       string
}

func NewAddFileCommand(
	logger logging.ILogger,
	storage clientEntities.IStorage,
	cryptoProvider clientEntities.ICryptoProvider,
	proxy clientEntities.IProxy,
	userID int,
) *AddFileCommand {
	return &AddFileCommand{
		Logger:         logger,
		Storage:        storage,
		CryptoProvider: cryptoProvider,
		proxy:          proxy,
		userID:         userID,
	}
}

func (c *AddFileCommand) GetName() string {
	return "add-file"
}

func (c *AddFileCommand) GetDescription() string {
	return "add binary file contents"
}

func (c *AddFileCommand) Validate(args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("example: add-file <path>")
	}
	c.filePath = args[0]
	return nil
}

func (c *AddFileCommand) Execute() cliEntities.CommandResult {
	meta, err := readNonSecretValue(c.Logger, "meta information")
	if err != nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("failed to read meta information: %v", err).Error(),
		}
	}
	rawData, err := readFile(c.filePath)
	if err != nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("failed to read file %s: %v", c.filePath, err).Error(),
		}
	}
	encrypted, err := c.CryptoProvider.Encrypt(rawData)
	if err != nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("failed to encrypt file data: %v", err).Error(),
		}
	}
	payload := common.FileReq{
		ServerID: nil,
		Meta:     meta,
		Value:    encrypted,
	}
	rowID, err := c.proxy.AddFile(&payload)
	if err != nil {
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
		ServerID: rowID,
	}
	if err := c.Storage.SaveFile(&fileLocal); err != nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("failed to store file data locally: %v", err).Error(),
		}
	}
	return cliEntities.CommandResult{
		SuccessMessage: "Successfully stored on server and locally",
	}
}

func readFile(filepath string) ([]byte, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return data, nil
}
