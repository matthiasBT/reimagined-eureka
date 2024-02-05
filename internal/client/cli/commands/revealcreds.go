package commands

import (
	"fmt"
	"strconv"

	cliEntities "reimagined_eureka/internal/client/cli/entities"
	clientEntities "reimagined_eureka/internal/client/entities"
	"reimagined_eureka/internal/client/infra/logging"
	"reimagined_eureka/internal/common"
)

type RevealCredsCommand struct {
	Logger         logging.ILogger
	Storage        clientEntities.IStorage
	CryptoProvider clientEntities.ICryptoProvider
	userID         int // TODO: check userID too!
	rowID          int
}

func NewRevealCredsCommand(
	logger logging.ILogger,
	storage clientEntities.IStorage,
	cryptoProvider clientEntities.ICryptoProvider,
	userID int,
) *RevealCredsCommand {
	return &RevealCredsCommand{
		Logger:         logger,
		Storage:        storage,
		CryptoProvider: cryptoProvider,
		userID:         userID,
	}
}

func (c *RevealCredsCommand) GetName() string {
	return "reveal-creds"
}

func (c *RevealCredsCommand) GetDescription() string {
	return "print the password of a login-password pair"
}

func (c *RevealCredsCommand) Validate(args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("example: reveal-creds <ID>")
	}
	rowID, err := strconv.Atoi(args[0])
	if err != nil || rowID <= 0 {
		return fmt.Errorf("value is a not a positive number")
	}
	c.rowID = rowID
	return nil
}

func (c *RevealCredsCommand) Execute() cliEntities.CommandResult {
	cred, err := c.Storage.ReadCredential(c.userID, c.rowID)
	if err != nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("failed to read password: %v", err).Error(),
		}
	}
	if cred == nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("credentials with ID %d don't exist for this user", c.rowID).Error(),
		}
	}
	encrypted := common.EncryptionResult{
		Ciphertext: cred.EncryptedPassword,
		Salt:       cred.Salt,
		Nonce:      cred.Nonce,
	}
	password, err := c.CryptoProvider.Decrypt(&encrypted)
	if err != nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("failed to decrypt password: %v", err).Error(),
		}
	}
	c.Logger.Warningln("Password:")
	c.Logger.Warningln(string(password))
	return cliEntities.CommandResult{}
}
