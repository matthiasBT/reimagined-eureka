package add

import (
	"fmt"

	"reimagined_eureka/internal/client/cli/commands"
	cliEntities "reimagined_eureka/internal/client/cli/entities"
	clientEntities "reimagined_eureka/internal/client/entities"
	"reimagined_eureka/internal/client/infra/logging"
	"reimagined_eureka/internal/common"
)

type AddCredsCommand struct {
	logger         logging.ILogger
	storage        clientEntities.IStorage
	cryptoProvider clientEntities.ICryptoProvider
	proxy          clientEntities.IProxy
	userID         int
	credsLogin     string
}

func NewAddCredsCommand(
	logger logging.ILogger,
	storage clientEntities.IStorage,
	cryptoProvider clientEntities.ICryptoProvider,
	proxy clientEntities.IProxy,
	userID int,
) *AddCredsCommand {
	return &AddCredsCommand{
		logger:         logger,
		storage:        storage,
		cryptoProvider: cryptoProvider,
		proxy:          proxy,
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
	encrypted, meta, err := PrepareCreds(c.logger, c.cryptoProvider)
	if err != nil {
		return cliEntities.CommandResult{
			FailureMessage: err.Error(),
		}
	}
	payload := common.CredentialsReq{
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
	if err := c.storage.SaveCredentials(&credsLocal); err != nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("failed to store credentials locally: %v", err).Error(),
		}
	}
	return cliEntities.CommandResult{
		SuccessMessage: "Successfully stored on server and locally",
	}
}

func PrepareCreds(
	logger logging.ILogger, cryptoProvider clientEntities.ICryptoProvider,
) (*common.EncryptionResult, string, error) {
	password, err := commands.ReadSecretValueMasked(logger, "password", 1, 0)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read password: %v", err)
	}
	meta, err := commands.ReadNonSecretValue(logger, "meta information")
	if err != nil {
		return nil, "", fmt.Errorf("failed to read meta information: %v", err)
	}
	encrypted, err := cryptoProvider.Encrypt([]byte(password))
	if err != nil {
		return nil, "", fmt.Errorf("failed to encrypt password: %v", err)
	}
	return encrypted, meta, nil
}
