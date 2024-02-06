package repositories

import (
	"context"
	"database/sql"
	"errors"

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

func (r *CredentialsRepo) Write(
	ctx context.Context, tx entities.Tx, userID int, data *common.CredentialsReq,
) (int, error) {
	if data.ServerID == nil {
		r.logger.Infof("Creating new credentials for user: %d", userID)
		return r.create(ctx, tx, userID, data)
	}
	r.logger.Infof("Updating credentials for user: %d", userID)
	return *data.ServerID, r.update(ctx, tx, userID, data)
}

func (r *CredentialsRepo) Read(
	ctx context.Context, tx entities.Tx, userID int, rowID int, lock bool,
) (*common.CredentialsReq, int, error) {
	var creds common.Credential
	query := "select * from credentials where id = $1 and user_id = $2 and not is_deleted"
	if lock {
		query = query + " for update"
	}
	if err := tx.GetContext(ctx, &creds, query, rowID, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, 0, entities.ErrDoesntExist
		}
		return nil, 0, err
	}
	var result common.CredentialsReq
	result.ServerID = &rowID
	result.Meta = creds.Meta
	result.Value = &common.EncryptionResult{
		Ciphertext: creds.EncryptedPassword,
		Salt:       creds.Salt,
		Nonce:      creds.Nonce,
	}
	return &result, creds.Version, nil
}

func (r *CredentialsRepo) ReadMany(
	ctx context.Context, tx entities.Tx, userID, startID, batchSize int,
) ([]*common.CredentialsReq, error) {
	var creds []common.Credential
	query := "select * from credentials where user_id = $1 and id > $2 and not is_deleted order by id limit $3"
	if err := tx.SelectContext(ctx, &creds, query, userID, startID, batchSize); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	var result = make([]*common.CredentialsReq, 0, len(creds))
	for _, row := range creds {
		resultRow := common.CredentialsReq{
			ServerID: &row.ID,
			Login:    row.Login,
			Meta:     row.Meta,
			Value: &common.EncryptionResult{
				Ciphertext: row.EncryptedPassword,
				Salt:       row.Salt,
				Nonce:      row.Nonce,
			},
		}
		result = append(result, &resultRow)
	}
	return result, nil
}

func (r *CredentialsRepo) Delete(ctx context.Context, tx entities.Tx, userID int, rowID int) error {
	_, _, err := r.Read(ctx, tx, userID, rowID, true)
	if err != nil {
		return err
	}
	query := "update credentials set is_deleted = true where id = $1"
	if err := tx.ExecContext(ctx, query, rowID); err != nil {
		r.logger.Errorf("Failed to delete creds: %s", err.Error())
		return err
	}
	return nil
}

func (r *CredentialsRepo) create(
	ctx context.Context, tx entities.Tx, userID int, data *common.CredentialsReq,
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
	r.logger.Infof("CredentialsReq created")
	if err := r.createVersion(ctx, tx, result.ID, data, entities.DefaultVersion); err != nil {
		return 0, err
	}
	return result.ID, nil
}

func (r *CredentialsRepo) createVersion(
	ctx context.Context, tx entities.Tx, credID int, data *common.CredentialsReq, version int,
) error {
	query := `
		insert into credentials_versions(cred_id, version, meta, login, encrypted_password, salt, nonce)
		values ($1, $2, $3, $4, $5, $6, $7)
	`
	if err := tx.ExecContext(
		ctx,
		query,
		credID,
		version,
		data.Meta,
		data.Login,
		data.Value.Ciphertext,
		data.Value.Salt,
		data.Value.Nonce,
	); err != nil {
		r.logger.Errorf("Failed to create credentials version: %s", err.Error())
		return err
	}
	r.logger.Infof("CredentialsReq version created")
	return nil
}

func (r *CredentialsRepo) update(ctx context.Context, tx entities.Tx, userID int, data *common.CredentialsReq) error {
	_, version, err := r.Read(ctx, tx, userID, *data.ServerID, true)
	if err != nil {
		return err
	}
	query := `
		update credentials
		set version = $2, meta = $3, encrypted_password = $4, salt = $5, nonce = $6, login = $7
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
		data.Login,
	); err != nil {
		r.logger.Errorf("Failed to update creds: %s", err.Error())
		return err
	}
	if err := r.createVersion(ctx, tx, *data.ServerID, data, version+1); err != nil {
		return err
	}
	return nil
}
