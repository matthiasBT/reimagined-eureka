package commands

import (
	"fmt"

	cliEntities "reimagined_eureka/internal/client/cli/entities"
	clientEntities "reimagined_eureka/internal/client/entities"
	"reimagined_eureka/internal/client/infra/logging"
	"reimagined_eureka/internal/common"
)

type AddCredsCommand struct {
	Logger         logging.ILogger
	Storage        clientEntities.IStorage
	CryptoProvider clientEntities.ICryptoProvider
	proxy          clientEntities.IProxy
	sessionCookie  string
	masterKey      string
	userID         int
	credsLogin     string
}

func NewAddCredentialsCommand(
	logger logging.ILogger,
	storage clientEntities.IStorage,
	cryptoProvider clientEntities.ICryptoProvider,
	proxy clientEntities.IProxy,
	sessionCookie string,
	masterKey string,
	userID int,
) *AddCredsCommand {
	return &AddCredsCommand{
		Logger:         logger,
		Storage:        storage,
		CryptoProvider: cryptoProvider,
		proxy:          proxy,
		sessionCookie:  sessionCookie,
		masterKey:      masterKey,
		userID:         userID,
	}
}

func (c *AddCredsCommand) GetName() string {
	return "add-creds"
}

func (c *AddCredsCommand) GetDescription() string {
	return "add a login-password pair"
}

func (c *AddCredsCommand) Validate(args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("example: add-creds <login>")
	}
	c.credsLogin = args[0]
	return nil
}

func (c *AddCredsCommand) Execute() cliEntities.CommandResult {
	password, err := readSecretValueMasked(c.Logger, "password", 1, 0)
	if err != nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("failed to read password: %v", err).Error(),
		}
	}
	meta, err := readNonSecretValue(c.Logger, "meta information")
	if err != nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("failed to read meta information: %v", err).Error(),
		}
	}
	encrypted, err := c.CryptoProvider.Encrypt([]byte(password))
	if err != nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("failed to encrypt password: %v", err).Error(),
		}
	}
	payload := common.Credentials{
		ServerID: nil,
		Login:    c.credsLogin,
		Meta:     meta,
		Value:    encrypted,
	}
	rowID, err := c.proxy.AddCredentials(&payload)
	if err != nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("request failed: %v", err).Error(),
		}
	}
	credsLocal := clientEntities.CredentialLocal{
		Credential: common.Credential{
			UserID:            c.userID,
			Meta:              payload.Meta,
			Login:             payload.Login,
			EncryptedPassword: payload.Value.Ciphertext,
			Salt:              payload.Value.Salt,
			Nonce:             payload.Value.Nonce,
		},
		ServerID: rowID,
	}
	if err := c.Storage.SaveCredentials(&credsLocal); err != nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("failed to store credentials locally: %v", err).Error(),
		}
	}
	return cliEntities.CommandResult{
		SuccessMessage: "Successfully stored on server and locally",
	}
}
