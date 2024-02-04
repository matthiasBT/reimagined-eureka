package commands

import (
	"fmt"

	cliEntities "reimagined_eureka/internal/client/cli/entities"
	clientEntities "reimagined_eureka/internal/client/entities"
	"reimagined_eureka/internal/client/infra/logging"
)

type ListSecretsCommand struct {
	Logger         logging.ILogger
	Storage        clientEntities.IStorage
	CryptoProvider clientEntities.ICryptoProvider
	userID         int // todo: check exported members of all structs!
	secretType     string
}

func NewListSecretsCommand(
	logger logging.ILogger,
	storage clientEntities.IStorage,
	cryptoProvider clientEntities.ICryptoProvider,
	userID int,
) *ListSecretsCommand {
	return &ListSecretsCommand{
		Logger:         logger,
		Storage:        storage,
		CryptoProvider: cryptoProvider,
		userID:         userID,
	}
}

func (c *ListSecretsCommand) GetName() string {
	return "list"
}

func (c *ListSecretsCommand) GetDescription() string {
	return "list all secrets or only those of some particular type"
}

func (c *ListSecretsCommand) Validate(args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("example: list <type>. Supported types: %s", listSupportedTypes())
	}
	secretType := args[0]
	if err := c.validateType(secretType); err != nil {
		return err
	}
	c.secretType = secretType
	return nil
}

func (c *ListSecretsCommand) Execute() cliEntities.CommandResult {
	var creds []*clientEntities.Credential
	var notes []*clientEntities.Note

	var err error
	if c.secretType == secretTypeCreds || c.secretType == secretTypeAll {
		if creds, err = c.Storage.ReadCredentials(c.userID); err != nil {
			return c.failedResult(err)
		}
	}
	if c.secretType == secretTypeNotes || c.secretType == secretTypeAll {
		if notes, err = c.Storage.ReadNotes(c.userID); err != nil {
			return c.failedResult(err)
		}
	}
	for _, cred := range creds {
		c.Logger.Warning("%d", cred.ID)
		c.Logger.Infoln("%s", cred.Purpose)
	}
	for _, note := range notes {
		c.Logger.Warning("%d", note.ID)
		c.Logger.Infoln("%s", note.Purpose)
	}
	msg := fmt.Sprintf(
		"Found %d rows. Show the decrypted contents of a secret using the \"reveal\" command",
		len(creds)+len(notes),
	)
	return cliEntities.CommandResult{SuccessMessage: msg}
}

func (c *ListSecretsCommand) validateType(what string) error {
	for _, tp := range supportedTypes {
		if what == tp {
			return nil
		}
	}
	return fmt.Errorf("unsupported secret type, must be one of those: %s", listSupportedTypes())
}

func (c *ListSecretsCommand) readCredentials() ([]*clientEntities.Credential, error) {
	return c.Storage.ReadCredentials(c.userID)
}

func (c *ListSecretsCommand) readNotes() ([]*clientEntities.Credential, error) {
	return c.Storage.ReadCredentials(c.userID)
}

func (c *ListSecretsCommand) failedResult(err error) cliEntities.CommandResult {
	return cliEntities.CommandResult{
		FailureMessage: fmt.Errorf("failed to list %s: %v", c.secretType, err).Error(),
	}
}
