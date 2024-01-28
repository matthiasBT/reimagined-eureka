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
	masterKey      string
}

func (c *MasterKeyCommand) GetName() string {
	return "master-key"
}

func (c *MasterKeyCommand) GetDescription() string {
	return "set (for new users) or check (for existing users) the master-key that will be used for encryption"
}

func (c *MasterKeyCommand) Validate(args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("example: master-key <login>")
	}
	//password, err := readSecretValueMasked(c.Logger, "user password")
	//if err != nil {
	//	return fmt.Errorf("failed to read user password: %v", err)
	//}
	//c.login, c.password = args[0], password
	return nil
}

func (c *MasterKeyCommand) Execute() cliEntities.CommandResult {
	panic("IMPLEMENT ME")
	//user, err := c.Storage.ReadUser(c.login)
	//if err != nil {
	//	return cliEntities.CommandResult{FailureMessage: err.Error()}
	//}
	//if user != nil {
	//	err := c.CryptoProvider.VerifyPassword(user, c.password)
	//	if err != nil {
	//		msg := fmt.Errorf("password verification failed: %v", err)
	//		return cliEntities.CommandResult{FailureMessage: msg.Error()}
	//	}
	//	return cliEntities.CommandResult{
	//		SuccessMessage: "Logged in successfully (locally)",
	//		// SessionCookie: nil  // TODO
	//		LoggedIn: true,
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
	//	LoggedIn:       true,
	//}
}
