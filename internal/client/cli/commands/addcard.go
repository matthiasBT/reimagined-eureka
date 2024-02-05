package commands

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	cliEntities "reimagined_eureka/internal/client/cli/entities"
	clientEntities "reimagined_eureka/internal/client/entities"
	"reimagined_eureka/internal/client/infra/logging"
	"reimagined_eureka/internal/common"
)

type AddCardCommand struct {
	Logger         logging.ILogger
	Storage        clientEntities.IStorage
	CryptoProvider clientEntities.ICryptoProvider
	proxy          clientEntities.IProxy
	sessionCookie  string
	masterKey      string
	userID         int
	cardNumber     string
}

func NewAddCardCommand(
	logger logging.ILogger,
	storage clientEntities.IStorage,
	cryptoProvider clientEntities.ICryptoProvider,
	proxy clientEntities.IProxy,
	sessionCookie string,
	masterKey string,
	userID int,
) *AddCardCommand {
	return &AddCardCommand{
		Logger:         logger,
		Storage:        storage,
		CryptoProvider: cryptoProvider,
		proxy:          proxy,
		sessionCookie:  sessionCookie,
		masterKey:      masterKey,
		userID:         userID,
	}
}

func (c *AddCardCommand) GetName() string {
	return "add-card"
}

func (c *AddCardCommand) GetDescription() string {
	return "add bank card data"
}

func (c *AddCardCommand) Validate(args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("example: add-card <number>")
	}
	c.cardNumber = args[0]
	return nil
}

func (c *AddCardCommand) Execute() cliEntities.CommandResult {
	monthRaw, err := readSecretValueMasked(c.Logger, "expiration date (month)", 1, 0)
	if err != nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("failed to read month: %v", err).Error(),
		}
	}
	yearRaw, err := readSecretValueMasked(c.Logger, "expiration date (year)", 1, 0)
	if err != nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("failed to read year: %v", err).Error(),
		}
	}
	cscRaw, err := readSecretValueMasked(c.Logger, "card security code", 1, 0)
	if err != nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("failed to read csc: %v", err).Error(),
		}
	}
	month, year, csc, err := parse(strings.Join([]string{monthRaw, yearRaw, cscRaw}, " "))
	if err != nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("failed to parse card data: %v", err).Error(),
		}
	}
	firstName, err := readSecretValueMasked(c.Logger, "owner (first name)", 1, 0)
	if err != nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("failed to read name: %v", err).Error(),
		}
	}
	lastName, err := readSecretValueMasked(c.Logger, "owner (last name)", 1, 0)
	if err != nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("failed to read name: %v", err).Error(),
		}
	}
	meta, err := readNonSecretValue(c.Logger, "meta information")
	if err != nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("failed to read meta information: %v", err).Error(),
		}
	}
	cardData := clientEntities.CardDataPlain{
		Month:     month,
		Year:      year,
		CSC:       csc,
		FirstName: firstName,
		LastName:  lastName,
	}
	cardBinary, err := json.Marshal(cardData)
	if err != nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("failed to prepare card data for encryption: %v", err).Error(),
		}
	}
	encrypted, err := c.CryptoProvider.Encrypt(cardBinary)
	if err != nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("failed to encrypt card data: %v", err).Error(),
		}
	}
	payload := common.CardReq{
		ServerID: nil,
		Meta:     meta,
		Value:    encrypted,
	}
	rowID, err := c.proxy.AddCard(&payload)
	if err != nil {
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
		ServerID: rowID,
	}
	if err := c.Storage.SaveCards(&cardLocal); err != nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("failed to store card data locally: %v", err).Error(),
		}
	}
	return cliEntities.CommandResult{
		SuccessMessage: "Successfully stored on server and locally",
	}
}

func parse(input string) (int, int, int, error) {
	re := regexp.MustCompile(`(?P<Month>\d{1,2})\s+(?P<Year>\d{2}|\d{4})\s+(?P<CSC>\d{3,4})`)
	matches := re.FindStringSubmatch(input)
	if matches == nil {
		return 0, 0, 0, fmt.Errorf("input does not match expected format")
	}

	matchMap := make(map[string]string)
	for i, name := range re.SubexpNames() {
		if i != 0 && name != "" {
			matchMap[name] = matches[i]
		}
	}

	month, err := strconv.Atoi(matchMap["Month"])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid month format")
	}

	year, err := strconv.Atoi(matchMap["Year"])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid year format")
	}

	csc, err := strconv.Atoi(matchMap["CSC"])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid csc format")
	}
	return month, year, csc, nil
}
