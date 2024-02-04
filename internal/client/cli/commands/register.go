package commands

import (
	"fmt"

	cliEntities "reimagined_eureka/internal/client/cli/entities"
	clientEntities "reimagined_eureka/internal/client/entities"
	"reimagined_eureka/internal/client/infra/logging"
	"reimagined_eureka/internal/common"
)

type RegisterCommand struct {
	LoginCommand
	masterKey, entropy string
}

func NewRegisterCommand(
	logger logging.ILogger,
	storage clientEntities.IStorage,
	proxy clientEntities.IProxy,
	cryptoProvider clientEntities.ICryptoProvider,
) *RegisterCommand {
	return &RegisterCommand{
		LoginCommand: LoginCommand{
			Logger:         logger,
			Storage:        storage,
			Proxy:          proxy,
			CryptoProvider: cryptoProvider,
		},
	}
}

func (c *RegisterCommand) GetName() string {
	return "register"
}

func (c *RegisterCommand) GetDescription() string {
	return "register a new user (requires a server connection)"
}

func (c *RegisterCommand) Validate(args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("example: register <login>")
	}
	login := args[0]
	if len(login) < common.MinLoginLength {
		return fmt.Errorf("login is shorter than %d characters", common.MinLoginLength)
	}
	password, err := readSecretValueMasked(c.Logger, "user password", common.MinPasswordLength, 0)
	if err != nil {
		return fmt.Errorf("failed to read user password: %v", err)
	}
	masterKey, err := readSecretValueMasked(c.Logger, "master key", minMasterKeyLength, maxMasterKeyLength)
	if err != nil {
		return fmt.Errorf("failed to read master key: %v", err)
	}
	c.Logger.Warningln("Now, it's time to create some random text that'll be used for master key verification")
	c.Logger.Warningln("You don't need to remember or store this text. Think of it as of entropy")
	entropy, err := readSecretValueMasked(c.Logger, "entropy", minEntropyLength, maxEntropyLength)
	if err != nil {
		return fmt.Errorf("failed to read entropy: %v", err)
	}
	c.LoginCommand.login = login
	c.LoginCommand.password = password
	c.masterKey = masterKey
	c.entropy = entropy
	return nil
}

func (c *RegisterCommand) Execute() cliEntities.CommandResult {
	tx, err := c.Storage.Tx()
	if err != nil {
		msg := fmt.Errorf("failed to register the user: %v", err)
		return cliEntities.CommandResult{FailureMessage: msg.Error()}
	}
	defer tx.Commit() // TODO: does it make sense to try to handle commit and rollback errors?
	c.CryptoProvider.SetMasterKey(c.masterKey)
	entropyBinary := []byte(c.entropy)
	entropyEncrypted, err := c.CryptoProvider.Encrypt(entropyBinary)
	if err != nil {
		msg := fmt.Errorf("failed to sign up: %v", err)
		defer tx.Rollback()
		return cliEntities.CommandResult{FailureMessage: msg.Error()}
	}
	entropyHash, err := c.CryptoProvider.Hash(entropyBinary)
	if err != nil {
		msg := fmt.Errorf("failed to sign up: %v", err)
		defer tx.Rollback()
		return cliEntities.CommandResult{FailureMessage: msg.Error()}
	}
	entropy := &common.Entropy{
		EncryptionResult: entropyEncrypted,
		Hash:             entropyHash,
	}
	userData, err := c.LoginCommand.Proxy.Register(c.LoginCommand.login, c.LoginCommand.password, entropy)
	if err != nil {
		msg := fmt.Errorf("failed to sign up: %v", err)
		defer tx.Rollback()
		return cliEntities.CommandResult{FailureMessage: msg.Error()}
	}
	newUser := &clientEntities.User{Login: c.LoginCommand.login}
	if err := c.LoginCommand.CryptoProvider.HashPassword(newUser, c.LoginCommand.password); err != nil {
		msg := fmt.Errorf("failed to store user %s data locally: %v", newUser.Login, err)
		defer tx.Rollback()
		return cliEntities.CommandResult{FailureMessage: msg.Error()}
	}
	userID, err := c.LoginCommand.Storage.SaveUser(newUser, entropy)
	if err != nil {
		msg := fmt.Errorf("failed to store user %s data locally: %v", newUser.Login, err)
		defer tx.Rollback()
		return cliEntities.CommandResult{FailureMessage: msg.Error()}
	}
	return cliEntities.CommandResult{
		SuccessMessage: "Registered successfully",
		SessionCookie:  userData.SessionCookie, // TODO: use it
		Login:          c.login,
		UserID:         userID,
	}
}
