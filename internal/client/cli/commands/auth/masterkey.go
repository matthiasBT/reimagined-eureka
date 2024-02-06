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
	Logger         logging.ILogger
	Storage        clientEntities.IStorage
	CryptoProvider clientEntities.ICryptoProvider
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
		Logger:         logger,
		Storage:        storage,
		CryptoProvider: cryptoProvider,
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
	key, err := commands.ReadSecretValueMasked(c.Logger, "master key", 0, 0) // TODO: fix 0s
	if err != nil {
		return fmt.Errorf("failed to read master key: %v", err)
	}
	c.masterKey = key
	c.CryptoProvider.SetMasterKey(key)
	return nil
}

func (c *MasterKeyCommand) Execute() cliEntities.CommandResult {
	c.Logger.Warningln("Checking the master key...")
	user, err := c.Storage.ReadUser(c.Login)
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
	entropy, err := c.CryptoProvider.Decrypt(&encrypted)
	if err != nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("failed to validate master-key: %v", err).Error(),
		}
	}
	if err := c.CryptoProvider.VerifyHash(entropy, user.EntropyHash); err != nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("incorrect master-key: %v", err).Error(),
		}
	}
	return cliEntities.CommandResult{
		SuccessMessage: "Master-key verified",
		MasterKey:      c.masterKey,
	}
}
