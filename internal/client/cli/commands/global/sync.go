package global

import (
	"fmt"

	"reimagined_eureka/internal/client/cli/commands"
	cliEntities "reimagined_eureka/internal/client/cli/entities"
	clientEntities "reimagined_eureka/internal/client/entities"
	"reimagined_eureka/internal/client/infra/logging"
	"reimagined_eureka/internal/common"
)

type SyncCommand struct {
	Logger         logging.ILogger
	Storage        clientEntities.IStorage
	CryptoProvider clientEntities.ICryptoProvider
	proxy          clientEntities.IProxy
	userID         int
}

func NewSyncCommand(
	logger logging.ILogger,
	storage clientEntities.IStorage,
	cryptoProvider clientEntities.ICryptoProvider,
	proxy clientEntities.IProxy,
	userID int,
) *SyncCommand {
	return &SyncCommand{
		Logger:         logger,
		Storage:        storage,
		CryptoProvider: cryptoProvider,
		proxy:          proxy,
		userID:         userID,
	}
}

func (c *SyncCommand) GetName() string {
	return "sync"
}

func (c *SyncCommand) GetDescription() string {
	return "retrieve a full data snapshot from the server and replace the local data"
}

func (c *SyncCommand) Validate(args ...string) error {
	if len(args) != 0 {
		return fmt.Errorf("example: sync")
	}
	return nil
}

func (c *SyncCommand) Execute() cliEntities.CommandResult {
	tx, err := c.Storage.Tx()
	if err != nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("failed to sync data: %v", err).Error(),
		}
	}
	defer tx.Commit()
	if err := c.Storage.Purge(c.userID); err != nil {
		defer tx.Rollback()
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("failed to clean old data: %v", err).Error(),
		}
	}
	if err := c.syncCredentials(); err != nil {
		defer tx.Rollback()
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("failed to sync credentials: %v", err).Error(),
		}
	}
	return cliEntities.CommandResult{
		SuccessMessage: "Sync completed",
	}
}

func (c *SyncCommand) syncCredentials() error {
	var startID int
	for {
		result, err := c.proxy.ReadCredentials(startID, commands.SyncCredsBatchSize)
		if err != nil {
			return err
		}
		if result == nil || len(result) == 0 {
			c.Logger.Warningln("Credentials synced")
			return nil
		}
		for _, row := range result {
			startID = *row.ServerID
			prepared := clientEntities.CredentialLocal{
				Credential: common.Credential{
					UserID:            c.userID,
					Meta:              row.Meta,
					Login:             row.Login,
					EncryptedPassword: row.Value.Ciphertext,
					Salt:              row.Value.Salt,
					Nonce:             row.Value.Nonce,
				},
				ServerID: *row.ServerID,
			}
			if err := c.Storage.SaveCredentials(&prepared); err != nil {
				return err
			}
		}
	}
}
