package add

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"reimagined_eureka/internal/client/cli/commands"
	cliEntities "reimagined_eureka/internal/client/cli/entities"
	clientEntities "reimagined_eureka/internal/client/entities"
	"reimagined_eureka/internal/client/infra/logging"
	"reimagined_eureka/internal/common"
)

type AddCardCommand struct {
	logger         logging.ILogger
	storage        clientEntities.IStorage
	cryptoProvider clientEntities.ICryptoProvider
	proxy          clientEntities.IProxy
	userID         int
	cardNumber     string
}

func NewAddCardCommand(
	logger logging.ILogger,
	storage clientEntities.IStorage,
	cryptoProvider clientEntities.ICryptoProvider,
	proxy clientEntities.IProxy,
	userID int,
) *AddCardCommand {
	return &AddCardCommand{
		logger:         logger,
		storage:        storage,
		cryptoProvider: cryptoProvider,
		proxy:          proxy,
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
		return fmt.Errorf("example: add-card <number> (without spaces)")
	}
	number := args[0]
	if !IsCardNumber(number) {
		return fmt.Errorf(
			"not a card number. Must contain only digits and be %d-%d digits long",
			commands.CardNumberMinLength,
			commands.CardNumberMaxLength,
		)
	}
	c.cardNumber = number
	return nil
}

func (c *AddCardCommand) Execute() cliEntities.CommandResult {
	encrypted, meta, err := PrepareCard(c.logger, c.cryptoProvider, c.cardNumber)
	if err != nil {
		return cliEntities.CommandResult{
			FailureMessage: err.Error(),
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
	if err := c.storage.SaveCard(&cardLocal); err != nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("failed to store card data locally: %v", err).Error(),
		}
	}
	return cliEntities.CommandResult{
		SuccessMessage: "Successfully stored on server and locally",
	}
}

func parse(input string) (string, string, string, error) {
	re := regexp.MustCompile(`(?P<Month>\d{1,2})\s+(?P<Year>\d{2}|\d{4})\s+(?P<CSC>\d{3,4})`)
	matches := re.FindStringSubmatch(input)
	if matches == nil {
		return "", "", "", fmt.Errorf("input does not match expected format")
	}

	matchMap := make(map[string]string)
	for i, name := range re.SubexpNames() {
		if i != 0 && name != "" {
			matchMap[name] = matches[i]
		}
	}

	month, err := strconv.Atoi(matchMap["Month"])
	if err != nil || month < commands.MonthMin || month > commands.MonthMax {
		return "", "", "", fmt.Errorf("invalid month format")
	}

	year, err := strconv.Atoi(matchMap["Year"])
	if err != nil || year < commands.YearMin {
		return "", "", "", fmt.Errorf("invalid year format")
	}

	_, err = strconv.Atoi(matchMap["CSC"])
	if err != nil {
		return "", "", "", fmt.Errorf("invalid csc format")
	}
	return matchMap["Month"], matchMap["Year"], matchMap["CSC"], nil
}

func IsCardNumber(what string) bool {
	for _, r := range what {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return !(len(what) < commands.CardNumberMinLength || len(what) > commands.CardNumberMaxLength)
}

func PrepareCard(
	logger logging.ILogger, cryptoProvider clientEntities.ICryptoProvider, cardNumber string,
) (*common.EncryptionResult, string, error) {
	monthRaw, err := commands.ReadSecretValueMasked(logger, "expiration date (month)", commands.MonthMinChars, commands.MonthMaxChars)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read month: %v", err)
	}
	yearRaw, err := commands.ReadSecretValueMasked(logger, "expiration date (year)", commands.YearMinChars, commands.YearMaxChars)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read year: %v", err)
	}
	cscRaw, err := commands.ReadSecretValueMasked(logger, "card security code", commands.CSCMinChars, commands.CSCMaxChars)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read csc: %v", err)
	}
	month, year, csc, err := parse(strings.Join([]string{monthRaw, yearRaw, cscRaw}, " "))
	if err != nil {
		return nil, "", fmt.Errorf("failed to parse card data: %v", err)
	}
	firstName, err := commands.ReadSecretValueMasked(logger, "owner (first name)", commands.NameMinChars, 0)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read name: %v", err)
	}
	lastName, err := commands.ReadSecretValueMasked(logger, "owner (last name)", commands.NameMinChars, 0)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read name: %v", err)
	}
	meta, err := commands.ReadNonSecretValue(logger, "meta information")
	if err != nil {
		return nil, "", fmt.Errorf("failed to read meta information: %v", err)
	}
	cardData := clientEntities.CardDataPlain{
		Month:     month,
		Year:      year,
		CSC:       csc,
		FirstName: firstName,
		LastName:  lastName,
		Number:    cardNumber,
	}
	cardBinary, err := json.Marshal(cardData)
	if err != nil {
		return nil, "", fmt.Errorf("failed to prepare card data for encryption: %v", err)
	}
	encrypted, err := cryptoProvider.Encrypt(cardBinary)
	if err != nil {
		return nil, "", fmt.Errorf("failed to encrypt card data: %v", err)
	}
	return encrypted, meta, nil
}
