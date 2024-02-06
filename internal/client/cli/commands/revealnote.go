package commands

import (
	"fmt"

	cliEntities "reimagined_eureka/internal/client/cli/entities"
	clientEntities "reimagined_eureka/internal/client/entities"
	"reimagined_eureka/internal/client/infra/logging"
	"reimagined_eureka/internal/common"
)

type RevealNoteCommand struct {
	logger         logging.ILogger
	storage        clientEntities.IStorage
	cryptoProvider clientEntities.ICryptoProvider
	userID         int
	rowID          int
	limit          int
}

func NewRevealNoteCommand(
	logger logging.ILogger,
	storage clientEntities.IStorage,
	cryptoProvider clientEntities.ICryptoProvider,
	userID int,
) *RevealNoteCommand {
	return &RevealNoteCommand{
		logger:         logger,
		storage:        storage,
		cryptoProvider: cryptoProvider,
		userID:         userID,
	}
}

func (c *RevealNoteCommand) GetName() string {
	return "reveal-note"
}

func (c *RevealNoteCommand) GetDescription() string {
	return "print the contents of a note"
}

func (c *RevealNoteCommand) Validate(args ...string) error {
	if len(args) < 1 || len(args) > 2 {
		return fmt.Errorf("example: reveal-note <ID> [<output-limit>]")
	}
	rowID, err := parsePositiveInt(args[0])
	if err != nil {
		return err
	}
	limit := 0
	if len(args) == 2 {
		limit, err = parsePositiveInt(args[1])
		if err != nil {
			return err
		}
	}
	c.rowID = rowID
	c.limit = limit
	return nil
}

func (c *RevealNoteCommand) Execute() cliEntities.CommandResult {
	note, err := c.storage.ReadNote(c.userID, c.rowID)
	if err != nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("failed to read note: %v", err).Error(),
		}
	}
	if note == nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("note %d doesn't exist for this user", c.rowID).Error(),
		}
	}
	encrypted := common.EncryptionResult{
		Ciphertext: note.EncryptedContent,
		Salt:       note.Salt,
		Nonce:      note.Nonce,
	}
	notePlain, err := c.cryptoProvider.Decrypt(&encrypted)
	if err != nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("failed to decrypt note: %v", err).Error(),
		}
	}
	c.logger.Warningln("Note:")
	text := string(notePlain)
	if c.limit != 0 {
		text = trimToNRunes(text, c.limit)
	}
	c.logger.Warningln(text)
	return cliEntities.CommandResult{}
}
