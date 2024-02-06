package commands

import (
	"encoding/json"
	"fmt"

	cliEntities "reimagined_eureka/internal/client/cli/entities"
	clientEntities "reimagined_eureka/internal/client/entities"
	"reimagined_eureka/internal/client/infra/logging"
	"reimagined_eureka/internal/common"
)

type RevealCardCommand struct {
	logger         logging.ILogger
	storage        clientEntities.IStorage
	cryptoProvider clientEntities.ICryptoProvider
	userID         int
	rowID          int
}

func NewRevealCardCommand(
	logger logging.ILogger,
	storage clientEntities.IStorage,
	cryptoProvider clientEntities.ICryptoProvider,
	userID int,
) *RevealCardCommand {
	return &RevealCardCommand{
		logger:         logger,
		storage:        storage,
		cryptoProvider: cryptoProvider,
		userID:         userID,
	}
}

func (c *RevealCardCommand) GetName() string {
	return "reveal-card"
}

func (c *RevealCardCommand) GetDescription() string {
	return "print the details of a bank card"
}

func (c *RevealCardCommand) Validate(args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("example: reveal-card <ID>")
	}
	rowID, err := parsePositiveInt(args[0])
	if err != nil {
		return err
	}
	c.rowID = rowID
	return nil
}

func (c *RevealCardCommand) Execute() cliEntities.CommandResult {
	card, err := c.storage.ReadCard(c.userID, c.rowID)
	if err != nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("failed to read card: %v", err).Error(),
		}
	}
	if card == nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("card with ID %d doesn't exist for this user", c.rowID).Error(),
		}
	}
	encrypted := common.EncryptionResult{
		Ciphertext: card.EncryptedContent,
		Salt:       card.Salt,
		Nonce:      card.Nonce,
	}
	cardBinary, err := c.cryptoProvider.Decrypt(&encrypted)
	if err != nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("failed to decrypt card: %v", err).Error(),
		}
	}
	var cardParsed clientEntities.CardDataPlain
	if err := json.Unmarshal(cardBinary, &cardParsed); err != nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("failed to parse card data: %v", err).Error(),
		}
	}
	c.logger.Warningln("Card details:")
	c.logger.Warningln("Number: %s", cardParsed.Number)
	c.logger.Warningln("Expiration date: %s/%s", cardParsed.Month, cardParsed.Year)
	c.logger.Warningln("Card security code (CSC): %s", cardParsed.CSC)
	c.logger.Warningln("Owner: %s %s", cardParsed.FirstName, cardParsed.LastName)
	return cliEntities.CommandResult{}
}
