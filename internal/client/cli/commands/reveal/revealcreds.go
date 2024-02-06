package reveal

import (
	"fmt"

	"reimagined_eureka/internal/client/cli/commands"
	cliEntities "reimagined_eureka/internal/client/cli/entities"
	clientEntities "reimagined_eureka/internal/client/entities"
	"reimagined_eureka/internal/client/infra/logging"
	"reimagined_eureka/internal/common"
)

type RevealCredsCommand struct {
	logger         logging.ILogger
	storage        clientEntities.IStorage
	cryptoProvider clientEntities.ICryptoProvider
	userID         int
	rowID          int
}

func NewRevealCredsCommand(
	logger logging.ILogger,
	storage clientEntities.IStorage,
	cryptoProvider clientEntities.ICryptoProvider,
	userID int,
) *RevealCredsCommand {
	return &RevealCredsCommand{
		logger:         logger,
		storage:        storage,
		cryptoProvider: cryptoProvider,
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
	rowID, err := commands.ParsePositiveInt(args[0])
	if err != nil {
		return err
	}
	c.rowID = rowID
	return nil
}

func (c *RevealCredsCommand) Execute() cliEntities.CommandResult {
	cred, err := c.storage.ReadCredential(c.userID, c.rowID)
	if err != nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("failed to read password: %v", err).Error(),
		}
	}
	if cred == nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("credentials %d don't exist", c.rowID).Error(),
		}
	}
	encrypted := common.EncryptionResult{
		Ciphertext: cred.EncryptedPassword,
		Salt:       cred.Salt,
		Nonce:      cred.Nonce,
	}
	password, err := c.cryptoProvider.Decrypt(&encrypted)
	if err != nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("failed to decrypt password: %v", err).Error(),
		}
	}
	c.logger.Warningln("Password:")
	c.logger.Warningln(string(password))
	return cliEntities.CommandResult{}
}
