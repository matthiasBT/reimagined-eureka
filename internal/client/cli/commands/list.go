package commands

import (
	"fmt"
	"strings"

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

// TODO: don't show deleted rows

func (c *ListSecretsCommand) Execute() cliEntities.CommandResult {
	if err := c.listCredentials(); err != nil {
		return c.failedResult(err)
	}
	if err := c.listNotes(); err != nil {
		return c.failedResult(err)
	}
	if err := c.listFiles(); err != nil {
		return c.failedResult(err)
	}
	if err := c.listCards(); err != nil {
		return c.failedResult(err)
	}
	msg := "Show the decrypted contents of a secret using the \"reveal\" command"
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

func (c *ListSecretsCommand) listCredentials() error {
	if c.secretType == secretTypeCreds || c.secretType == secretTypeAll {
		c.Logger.Warningln("Credentials in storage:")
		if creds, err := c.Storage.ReadCredentials(c.userID); err != nil {
			return err
		} else {
			for _, cred := range creds {
				c.printItem(cred.ID, cred.Meta, cred.Login)
			}
		}
	}
	return nil
}

func (c *ListSecretsCommand) listNotes() error {
	if c.secretType == secretTypeNotes || c.secretType == secretTypeAll {
		c.Logger.Warningln("Notes in storage:")
		if notes, err := c.Storage.ReadNotes(c.userID); err != nil {
			return err
		} else {
			for _, note := range notes {
				c.printItem(note.ID, note.Meta)
			}
		}
	}
	return nil
}

func (c *ListSecretsCommand) listFiles() error {
	if c.secretType == secretTypeFiles || c.secretType == secretTypeAll {
		c.Logger.Warningln("Files in storage:")
		if files, err := c.Storage.ReadFiles(c.userID); err != nil {
			return err
		} else {
			for _, file := range files {
				c.printItem(file.ID, file.Meta)
			}
		}
	}
	return nil
}

func (c *ListSecretsCommand) listCards() error {
	if c.secretType == secretTypeCards || c.secretType == secretTypeAll {
		c.Logger.Warningln("Cards in storage:")
		if cards, err := c.Storage.ReadCards(c.userID); err != nil {
			return err
		} else {
			for _, card := range cards {
				c.printItem(card.ID, card.Meta)
			}
		}
	}
	return nil
}

func (c *ListSecretsCommand) failedResult(err error) cliEntities.CommandResult {
	return cliEntities.CommandResult{
		FailureMessage: fmt.Errorf("failed to list %s: %v. Aborting", c.secretType, err).Error(),
	}
}

func (c *ListSecretsCommand) printItem(id int, args ...string) {
	c.Logger.Warningln("ID: %d", id)
	for _, arg := range args {
		c.Logger.Infoln("%s", arg)
	}
	c.Logger.Infoln(strings.Repeat(secretDelimiterChar, secretDelimiterWidth))
}
