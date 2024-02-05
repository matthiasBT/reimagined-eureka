package repositories

import (
	"context"

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

func (r *CardsRepo) Read(ctx context.Context, tx entities.Tx, userID int, rowId int) (*common.CardReq, error) {
	panic("implement me!")
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
	if err := r.createVersion(ctx, tx, result.ID, data); err != nil {
		return 0, err
	}
	return result.ID, nil
}

func (r *CardsRepo) createVersion(
	ctx context.Context, tx entities.Tx, cardID int, data *common.CardReq,
) error {
	query := `
		insert into cards_versions(card_id, version, meta, encrypted_content, salt, nonce)
		values ($1, $2, $3, $4, $5, $6)
	`
	if err := tx.ExecContext(
		ctx,
		query,
		cardID,
		entities.DefaultVersion,
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

func (r *CardsRepo) update(ctx context.Context, tx entities.Tx, userID int, data *common.CardReq) (int, error) {
	panic("Implement me!")
}
