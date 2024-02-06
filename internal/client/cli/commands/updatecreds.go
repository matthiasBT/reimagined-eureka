package commands

import (
	"fmt"

	cliEntities "reimagined_eureka/internal/client/cli/entities"
	clientEntities "reimagined_eureka/internal/client/entities"
	"reimagined_eureka/internal/client/infra/logging"
	"reimagined_eureka/internal/common"
)

type UpdateCredsCommand struct {
	Logger         logging.ILogger
	Storage        clientEntities.IStorage
	CryptoProvider clientEntities.ICryptoProvider
	proxy          clientEntities.IProxy
	userID         int
	rowID          int
	credsLogin     string
}

func NewUpdateCredsCommand(
	logger logging.ILogger,
	storage clientEntities.IStorage,
	cryptoProvider clientEntities.ICryptoProvider,
	proxy clientEntities.IProxy,
	userID int,
) *UpdateCredsCommand {
	return &UpdateCredsCommand{
		Logger:         logger,
		Storage:        storage,
		CryptoProvider: cryptoProvider,
		proxy:          proxy,
		userID:         userID,
	}
}

func (c *UpdateCredsCommand) GetName() string {
	return "update-creds"
}

func (c *UpdateCredsCommand) GetDescription() string {
	return "update a login-password pair"
}

func (c *UpdateCredsCommand) Validate(args ...string) error {
	if len(args) != 2 {
		return fmt.Errorf("example: update-creds <ID> <login>")
	}
	rowID, err := parsePositiveInt(args[0])
	if err != nil {
		return err
	}
	creds, err := c.Storage.ReadCredential(c.userID, rowID)
	if err != nil {
		return fmt.Errorf("failed to read creds: %v", err)
	}
	if creds == nil {
		return fmt.Errorf("creds %d don't exist for this user", rowID)
	}
	c.rowID = rowID
	c.credsLogin = args[1]
	return nil
}

func (c *UpdateCredsCommand) Execute() cliEntities.CommandResult {
	encrypted, meta, err := prepareCreds(c.Logger, c.CryptoProvider)
	if err != nil {
		return cliEntities.CommandResult{
			FailureMessage: err.Error(),
		}
	}
	payload := common.CredentialsReq{
		ServerID: &c.rowID,
		Login:    c.credsLogin,
		Meta:     meta,
		Value:    encrypted,
	}
	if err := c.proxy.UpdateCredentials(&payload); err != nil {
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
		ServerID: c.rowID,
	}
	if err := c.Storage.SaveCredentials(&credsLocal); err != nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("failed to update credentials locally: %v", err).Error(),
		}
	}
	return cliEntities.CommandResult{
		SuccessMessage: "Successfully updated on server and locally",
	}
}
