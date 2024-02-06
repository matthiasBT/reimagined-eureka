package repositories

import (
	"context"
	"database/sql"
	"errors"

	"reimagined_eureka/internal/common"
	"reimagined_eureka/internal/server/entities"
	"reimagined_eureka/internal/server/infra/logging"
)

type CardsRepo struct {
	logger  logging.ILogger
	storage entities.Storage
}

func NewCardsRepo(logger logging.ILogger, storage entities.Storage) *CardsRepo {
	return &CardsRepo{
		logger:  logger,
		storage: storage,
	}
}

func (r *CardsRepo) Write(ctx context.Context, tx entities.Tx, userID int, data *common.CardReq) (int, error) {
	r.logger.Infof("Creating new card for user: %d", userID)
	return r.create(ctx, tx, userID, data)
}

func (r *CardsRepo) Read(
	ctx context.Context, tx entities.Tx, userID int, rowID int, lock bool,
) (*common.CardReq, int, error) {
	var card common.Card
	query := "select * from cards where id = $1 and user_id = $2" // TODO: check delete flag in the future
	if lock {
		query = query + " for update"
	}
	if err := tx.GetContext(ctx, &card, query, rowID, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, 0, entities.ErrDoesntExist
		}
		return nil, 0, err
	}
	var result common.CardReq
	result.ServerID = &rowID
	result.Meta = card.Meta
	result.Value = &common.EncryptionResult{
		Ciphertext: card.EncryptedContent,
		Salt:       card.Salt,
		Nonce:      card.Nonce,
	}
	return &result, card.Version, nil
}

func (r *CardsRepo) create(
	ctx context.Context, tx entities.Tx, userID int, data *common.CardReq,
) (int, error) {
	var result common.Card
	query := `
		insert into cards(user_id, meta, encrypted_content, salt, nonce)
		values ($1, $2, $3, $4, $5)
		returning *
	`
	if err := tx.GetContext(
		ctx, &result, query, userID, data.Meta, data.Value.Ciphertext, data.Value.Salt, data.Value.Nonce,
	); err != nil {
		r.logger.Errorf("Failed to create card: %s", err.Error())
		return 0, err
	}
	r.logger.Infof("Card created")
	if err := r.createVersion(ctx, tx, result.ID, data, entities.DefaultVersion); err != nil {
		return 0, err
	}
	return result.ID, nil
}

func (r *CardsRepo) createVersion(
	ctx context.Context, tx entities.Tx, cardID int, data *common.CardReq, version int,
) error {
	query := `
		insert into cards_versions(card_id, version, meta, encrypted_content, salt, nonce)
		values ($1, $2, $3, $4, $5, $6)
	`
	if err := tx.ExecContext(
		ctx,
		query,
		cardID,
		version,
		data.Meta,
		data.Value.Ciphertext,
		data.Value.Salt,
		data.Value.Nonce,
	); err != nil {
		r.logger.Errorf("Failed to create card version: %s", err.Error())
		return err
	}
	r.logger.Infof("Card version created")
	return nil
}

func (r *CardsRepo) update(ctx context.Context, tx entities.Tx, userID int, data *common.CardReq) error {
	_, version, err := r.Read(ctx, tx, userID, *data.ServerID, true)
	if err != nil {
		return err
	}
	query := `
		update cards
		set version = $2, meta = $3, encrypted_content = $4, salt = $5, nonce = $6
		where id = $1
	`
	if err := tx.ExecContext(
		ctx,
		query,
		*data.ServerID,
		version+1,
		data.Meta,
		data.Value.Ciphertext,
		data.Value.Salt,
		data.Value.Nonce,
	); err != nil {
		r.logger.Errorf("Failed to update card: %s", err.Error())
		return err
	}
	if err := r.createVersion(ctx, tx, *data.ServerID, data, version+1); err != nil {
		return err
	}
	return nil
}
