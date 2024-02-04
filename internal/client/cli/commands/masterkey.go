package commands

import (
	"fmt"

	cliEntities "reimagined_eureka/internal/client/cli/entities"
	clientEntities "reimagined_eureka/internal/client/entities"
	"reimagined_eureka/internal/client/infra/logging"
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
	return "set (for new users) or check (for existing users) the master-key that will be used for encryption"
}

func (c *MasterKeyCommand) Validate(args ...string) error {
	if len(args) != 0 {
		return fmt.Errorf("example: master-key")
	}
	key, err := readSecretValueMasked(c.Logger, "master key", 0, 0)
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
	entropy, err := c.CryptoProvider.Decrypt(user.EntropyEncrypted, user.EntropySalt, user.EntropyNonce)
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
		// SessionCookie: nil  // TODO: pass along?
		MasterKeyVerified: true,
	}

	//if user != nil {
	//	err := c.CryptoProvider.VerifyPassword(user, c.password)
	//	if err != nil {
	//		msg := fmt.Errorf("password verification failed: %v", err)
	//		return cliEntities.CommandResult{FailureMessage: msg.Error()}
	//	}
	//	return cliEntities.CommandResult{
	//		SuccessMessage: "Logged in successfully (locally)",
	//		// SessionCookie: nil  // TODO
	//		Login: true,
	//	}
	//}
	//c.Logger.Warningln("User %s not found locally. Going to fetch it from server", c.login)
	//userData, err := c.Proxy.LogIn(c.login, c.password)
	//if err != nil {
	//	msg := fmt.Errorf("failed to log in: %v", err)
	//	return cliEntities.CommandResult{FailureMessage: msg.Error()}
	//}
	//newUser := &clientEntities.User{Login: c.login}
	//if err := c.CryptoProvider.HashPassword(newUser, c.password); err != nil {
	//	msg := fmt.Errorf("failed to store user %s data locally: %v", newUser.Login, err)
	//	return cliEntities.CommandResult{FailureMessage: msg.Error()}
	//}
	//if err := c.Storage.SaveUser(newUser); err != nil {
	//	msg := fmt.Errorf("failed to store user %s data locally: %v", newUser.Login, err)
	//	return cliEntities.CommandResult{FailureMessage: msg.Error()}
	//}
	//return cliEntities.CommandResult{
	//	SuccessMessage: "Logged in successfully (on server)",
	//	SessionCookie:  userData.SessionCookie,
	//	Login:       true,
	//}
}
