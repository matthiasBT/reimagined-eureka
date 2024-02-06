package commands

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	cliEntities "reimagined_eureka/internal/client/cli/entities"
	clientEntities "reimagined_eureka/internal/client/entities"
	"reimagined_eureka/internal/client/infra/logging"
	"reimagined_eureka/internal/common"
)

const monthMin = 1
const monthMax = 12
const yearMin = 20
const monthMinChars = 1
const monthMaxChars = 1
const yearMinChars = 1
const yearMaxChars = 2
const cscMinChars = 3
const cscMaxChars = 3
const nameMinChars = 1

type AddCardCommand struct {
	Logger         logging.ILogger
	Storage        clientEntities.IStorage
	CryptoProvider clientEntities.ICryptoProvider
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
		Logger:         logger,
		Storage:        storage,
		CryptoProvider: cryptoProvider,
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
	if !isCardNumber(number) {
		return fmt.Errorf(
			"not a card number. Must contain only digits and be %d-%d digits long",
			cardNumberMinLength,
			cardNumberMaxLength,
		)
	}
	c.cardNumber = number
	return nil
}

func (c *AddCardCommand) Execute() cliEntities.CommandResult {
	monthRaw, err := readSecretValueMasked(c.Logger, "expiration date (month)", monthMinChars, monthMaxChars)
	if err != nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("failed to read month: %v", err).Error(),
		}
	}
	yearRaw, err := readSecretValueMasked(c.Logger, "expiration date (year)", yearMinChars, yearMaxChars)
	if err != nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("failed to read year: %v", err).Error(),
		}
	}
	cscRaw, err := readSecretValueMasked(c.Logger, "card security code", cscMinChars, cscMaxChars)
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
	firstName, err := readSecretValueMasked(c.Logger, "owner (first name)", nameMinChars, 0)
	if err != nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("failed to read name: %v", err).Error(),
		}
	}
	lastName, err := readSecretValueMasked(c.Logger, "owner (last name)", nameMinChars, 0)
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
		Number:    c.cardNumber,
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
	if err != nil || month < monthMin || month > monthMax {
		return "", "", "", fmt.Errorf("invalid month format")
	}

	year, err := strconv.Atoi(matchMap["Year"])
	if err != nil || year < yearMin {
		return "", "", "", fmt.Errorf("invalid year format")
	}

	_, err = strconv.Atoi(matchMap["CSC"])
	if err != nil {
		return "", "", "", fmt.Errorf("invalid csc format")
	}
	return matchMap["Month"], matchMap["Year"], matchMap["CSC"], nil
}

func isCardNumber(what string) bool {
	for _, r := range what {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return !(len(what) < cardNumberMinLength || len(what) > cardNumberMaxLength)
}
