package repositories

import (
	"context"

	"reimagined_eureka/internal/common"
	"reimagined_eureka/internal/server/entities"
	"reimagined_eureka/internal/server/infra/logging"
)

type CredentialsRepo struct {
	logger  logging.ILogger
	storage entities.Storage
}

func NewCredentialsRepo(logger logging.ILogger, storage entities.Storage) *CredentialsRepo {
	return &CredentialsRepo{
		logger:  logger,
		storage: storage,
	}
}

func (r *CredentialsRepo) Write(ctx context.Context, tx entities.Tx, userID int, data *common.Credentials) (int, error) {
	r.logger.Infof("Creating new credentials for user: %d", userID)
	return r.create(ctx, tx, userID, data)
}

func (r *CredentialsRepo) Read(ctx context.Context, tx entities.Tx, userID int, rowId int) (*common.Credentials, error) {
	panic("implement me!")
}

func (r *CredentialsRepo) create(
	ctx context.Context, tx entities.Tx, userID int, data *common.Credentials,
) (int, error) {
	var result common.Credential
	query := `
		insert into credentials(user_id, meta, login, encrypted_password, salt, nonce)
		values ($1, $2, $3, $4, $5, $6)
		returning *
	`
	if err := tx.GetContext(
		ctx, &result, query, userID, data.Meta, data.Login, data.Value.Ciphertext, data.Value.Salt, data.Value.Nonce,
	); err != nil {
		r.logger.Errorf("Failed to create credentials: %s", err.Error())
		return 0, err
	}
	r.logger.Infof("Credentials created")
	if err := r.createVersion(ctx, tx, result.ID, data); err != nil {
		return 0, err
	}
	return result.ID, nil
}

func (r *CredentialsRepo) createVersion(
	ctx context.Context, tx entities.Tx, credID int, data *common.Credentials,
) error {
	query := `
		insert into credentials_versions(cred_id, version, meta, login, encrypted_password, salt, nonce)
		values ($1, $2, $3, $4, $5, $6, $7)
	`
	if err := tx.ExecContext(
		ctx,
		query,
		credID,
		entities.DefaultVersion,
		data.Meta,
		data.Login,
		data.Value.Ciphertext,
		data.Value.Salt,
		data.Value.Nonce,
	); err != nil {
		r.logger.Errorf("Failed to create credentials version: %s", err.Error())
		return err
	}
	r.logger.Infof("Credentials version created")
	return nil
}

func (r *CredentialsRepo) update(ctx context.Context, tx entities.Tx, userID int, data *common.Credentials) (int, error) {
	panic("Implement me!")
}
