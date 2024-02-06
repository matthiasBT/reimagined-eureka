package auth

import (
	"fmt"

	"reimagined_eureka/internal/client/cli/commands"
	cliEntities "reimagined_eureka/internal/client/cli/entities"
	clientEntities "reimagined_eureka/internal/client/entities"
	"reimagined_eureka/internal/client/infra/logging"
	"reimagined_eureka/internal/common"
)

type MasterKeyCommand struct {
	logger         logging.ILogger
	storage        clientEntities.IStorage
	cryptoProvider clientEntities.ICryptoProvider
	Login          string
	masterKey      string
}

func NewMasterKeyCommand(
	logger logging.ILogger,
	storage clientEntities.IStorage,
	cryptoProvider clientEntities.ICryptoProvider,
	login string,
) *MasterKeyCommand {
	return &MasterKeyCommand{
		logger:         logger,
		storage:        storage,
		cryptoProvider: cryptoProvider,
		Login:          login,
	}
}

func (c *MasterKeyCommand) GetName() string {
	return "master-key"
}

func (c *MasterKeyCommand) GetDescription() string {
	return "check the master-key that will be used for encryption"
}

func (c *MasterKeyCommand) Validate(args ...string) error {
	if len(args) != 0 {
		return fmt.Errorf("example: master-key")
	}
	key, err := commands.ReadSecretValueMasked(
		c.logger, "master key", commands.MinMasterKeyLength, commands.MaxMasterKeyLength,
	)
	if err != nil {
		return fmt.Errorf("failed to read master key: %v", err)
	}
	c.masterKey = key
	c.cryptoProvider.SetMasterKey(key)
	return nil
}

func (c *MasterKeyCommand) Execute() cliEntities.CommandResult {
	c.logger.Warningln("Checking the master key...")
	user, err := c.storage.ReadUser(c.Login)
	if err != nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("failed to validate master-key: %v", err).Error(),
		}
	}
	encrypted := common.EncryptionResult{
		Ciphertext: user.EntropyEncrypted,
		Salt:       user.EntropySalt,
		Nonce:      user.EntropyNonce,
	}
	entropy, err := c.cryptoProvider.Decrypt(&encrypted)
	if err != nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("failed to validate master-key: %v", err).Error(),
		}
	}
	if err := c.cryptoProvider.VerifyHash(entropy, user.EntropyHash); err != nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("incorrect master-key: %v", err).Error(),
		}
	}
	return cliEntities.CommandResult{
		SuccessMessage: "Master-key verified",
		MasterKey:      c.masterKey,
	}
}
