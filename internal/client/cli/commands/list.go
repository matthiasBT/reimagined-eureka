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
	Login          string // todo: check exported members of all structs!
	secretType     string
}

func NewListSecretsCommand(
	logger logging.ILogger,
	storage clientEntities.IStorage,
	cryptoProvider clientEntities.ICryptoProvider,
	login string,
) *ListSecretsCommand {
	return &ListSecretsCommand{
		Logger:         logger,
		Storage:        storage,
		CryptoProvider: cryptoProvider,
		Login:          login,
	}
}

func (c *ListSecretsCommand) GetName() string {
	return "list"
}

func (c *ListSecretsCommand) GetDescription() string {
	return "list all secrets or only those of some particular type"
}

func (c *ListSecretsCommand) Validate(args ...string) error {
	if len(args) > 1 {
		return fmt.Errorf("example: list [<type>]")
	}
	if len(args) == 0 {
		c.secretType = ""
	} else if err := c.validateType(args[0]); err != nil {
		return err
	} else {
		c.secretType = args[0]
	}
	return nil
}

func (c *ListSecretsCommand) Execute() cliEntities.CommandResult {
	result, err := c.Storage.ReadCredentials(c.Login, c.secretType)
	if err != nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("failed to list %s master-key: %v", c.secretType, err).Error(),
		}
	}
	for _, creds := range result {
		c.Logger.Warning("%d", creds.ID)
		c.Logger.Infoln("%s", creds.Purpose)
	}
	msg := fmt.Sprintf(
		"Found %d rows. Show the decrypted contents of a secret using the \"reveal\" command", len(result),
	)
	return cliEntities.CommandResult{SuccessMessage: msg}
}

func (c *ListSecretsCommand) validateType(what string) error {
	for _, tp := range supportedTypes {
		if what == tp {
			return nil
		}
	}
	hrTypes := strings.TrimSpace(strings.Join(supportedTypes, ", "))
	return fmt.Errorf("unsupported secret type, must be one of those: %s", hrTypes)
}
