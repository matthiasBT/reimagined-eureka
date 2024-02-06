package commands

import (
	"bufio"
	"fmt"
	"os"

	cliEntities "reimagined_eureka/internal/client/cli/entities"
	clientEntities "reimagined_eureka/internal/client/entities"
	"reimagined_eureka/internal/client/infra/logging"
	"reimagined_eureka/internal/common"
)

type RevealFileCommand struct {
	logger         logging.ILogger
	storage        clientEntities.IStorage
	cryptoProvider clientEntities.ICryptoProvider
	userID         int
	rowID          int
	savePath       string
	limit          int
}

func NewRevealFileCommand(
	logger logging.ILogger,
	storage clientEntities.IStorage,
	cryptoProvider clientEntities.ICryptoProvider,
	userID int,
) *RevealFileCommand {
	return &RevealFileCommand{
		logger:         logger,
		storage:        storage,
		cryptoProvider: cryptoProvider,
		userID:         userID,
	}
}

func (c *RevealFileCommand) GetName() string {
	return "reveal-file"
}

func (c *RevealFileCommand) GetDescription() string {
	return "write the contents of a file to the destination path"
}

func (c *RevealFileCommand) Validate(args ...string) error {
	if len(args) != 2 {
		return fmt.Errorf("example: reveal-file <ID> <local path>")
	}
	rowID, err := parsePositiveInt(args[0])
	if err != nil {
		return err
	}
	c.rowID = rowID
	c.savePath = args[1]
	return nil
}

func (c *RevealFileCommand) Execute() cliEntities.CommandResult {
	file, err := c.storage.ReadFile(c.userID, c.rowID)
	if err != nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("failed to read file: %v", err).Error(),
		}
	}
	if file == nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("file %d doesn't exist for this user", c.rowID).Error(),
		}
	}
	encrypted := common.EncryptionResult{
		Ciphertext: file.EncryptedContent,
		Salt:       file.Salt,
		Nonce:      file.Nonce,
	}
	filePlain, err := c.cryptoProvider.Decrypt(&encrypted)
	if err != nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("failed to decrypt file: %v", err).Error(),
		}
	}
	c.logger.Warningln("Saving the file to %s...", c.savePath)
	if err := dumpFile(filePlain, c.savePath); err != nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("failed to save file: %v", err).Error(),
		}
	}
	return cliEntities.CommandResult{
		SuccessMessage: "File saved!",
	}
}

func dumpFile(what []byte, where string) error {
	file, err := os.Create(where)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	_, err = writer.Write(what)
	if err != nil {
		return err
	}
	err = writer.Flush()
	if err != nil {
		return err
	}
	return nil
}
