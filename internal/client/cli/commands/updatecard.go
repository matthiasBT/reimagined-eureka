package commands

import (
	"fmt"

	cliEntities "reimagined_eureka/internal/client/cli/entities"
	clientEntities "reimagined_eureka/internal/client/entities"
	"reimagined_eureka/internal/client/infra/logging"
	"reimagined_eureka/internal/common"
)

type UpdateCardCommand struct {
	Logger         logging.ILogger
	Storage        clientEntities.IStorage
	CryptoProvider clientEntities.ICryptoProvider
	proxy          clientEntities.IProxy
	userID         int
	rowIDServer    int
	rowIDLocal     int
	cardNumber     string
}

func NewUpdateCardCommand(
	logger logging.ILogger,
	storage clientEntities.IStorage,
	cryptoProvider clientEntities.ICryptoProvider,
	proxy clientEntities.IProxy,
	userID int,
) *UpdateCardCommand {
	return &UpdateCardCommand{
		Logger:         logger,
		Storage:        storage,
		CryptoProvider: cryptoProvider,
		proxy:          proxy,
		userID:         userID,
	}
}

func (c *UpdateCardCommand) GetName() string {
	return "update-card"
}

func (c *UpdateCardCommand) GetDescription() string {
	return "update bank card data"
}

func (c *UpdateCardCommand) Validate(args ...string) error {
	if len(args) != 2 {
		return fmt.Errorf("example: update-card <ID> <number> (without spaces)")
	}
	rowID, err := parsePositiveInt(args[0])
	if err != nil {
		return err
	}
	card, err := c.Storage.ReadCard(c.userID, rowID)
	if err != nil {
		return fmt.Errorf("failed to read card: %v", err)
	}
	if card == nil {
		return fmt.Errorf("card %d doesn't exist", rowID)
	}
	number := args[1]
	if !isCardNumber(number) {
		return fmt.Errorf(
			"not a card number. Must contain only digits and be %d-%d digits long",
			cardNumberMinLength,
			cardNumberMaxLength,
		)
	}
	c.rowIDServer = card.ServerID
	c.rowIDLocal = card.ID
	c.cardNumber = number
	return nil
}

func (c *UpdateCardCommand) Execute() cliEntities.CommandResult {
	encrypted, meta, err := prepareCard(c.Logger, c.CryptoProvider, c.cardNumber)
	if err != nil {
		return cliEntities.CommandResult{
			FailureMessage: err.Error(),
		}
	}
	payload := common.CardReq{
		ServerID: &c.rowIDServer,
		Meta:     meta,
		Value:    encrypted,
	}
	if err := c.proxy.UpdateCard(&payload); err != nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("request failed: %v", err).Error(),
		}
	}
	cardLocal := clientEntities.CardLocal{
		Card: common.Card{
			UserID:           c.userID,
			Meta:             payload.Meta,
			EncryptedContent: payload.Value.Ciphertext,
			Salt:             payload.Value.Salt,
			Nonce:            payload.Value.Nonce,
		},
		ServerID: c.rowIDServer,
	}
	if err := c.Storage.SaveCard(&cardLocal); err != nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("failed to update card data locally: %v", err).Error(),
		}
	}
	return cliEntities.CommandResult{
		SuccessMessage: "Successfully updated on server and locally",
	}
}
