package repositories

import (
	"context"
	"database/sql"
	"errors"

	"reimagined_eureka/internal/common"
	"reimagined_eureka/internal/server/entities"
	"reimagined_eureka/internal/server/infra/logging"
)

type FilesRepo struct {
	logger  logging.ILogger
	storage entities.Storage
}

func NewFilesRepo(logger logging.ILogger, storage entities.Storage) *FilesRepo {
	return &FilesRepo{
		logger:  logger,
		storage: storage,
	}
}

func (r *FilesRepo) Write(ctx context.Context, tx entities.Tx, userID int, data *common.FileReq) (int, error) {
	if data.ServerID == nil {
		r.logger.Infof("Creating new file for user: %d", userID)
		return r.create(ctx, tx, userID, data)
	}
	r.logger.Infof("Updating file %d for user: %d", data.ServerID, userID)
	return *data.ServerID, r.update(ctx, tx, userID, data)
}

func (r *FilesRepo) Read(
	ctx context.Context, tx entities.Tx, userID int, rowID int, lock bool,
) (*common.FileReq, int, error) {
	var file common.File
	query := "select * from files where id = $1 and user_id = $2 and not is_deleted"
	if lock {
		query = query + " for update"
	}
	if err := tx.GetContext(ctx, &file, query, rowID, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, 0, entities.ErrDoesntExist
		}
		return nil, 0, err
	}
	var result common.FileReq
	result.ServerID = &rowID
	result.Meta = file.Meta
	result.Value = &common.EncryptionResult{
		Ciphertext: file.EncryptedContent,
		Salt:       file.Salt,
		Nonce:      file.Nonce,
	}
	return &result, file.Version, nil
}

func (r *FilesRepo) Delete(ctx context.Context, tx entities.Tx, userID int, rowID int) error {
	_, _, err := r.Read(ctx, tx, userID, rowID, true)
	if err != nil {
		return err
	}
	query := "update files set is_deleted = true where id = $1"
	if err := tx.ExecContext(ctx, query, rowID); err != nil {
		r.logger.Errorf("Failed to delete file: %s", err.Error())
		return err
	}
	return nil
}

func (r *FilesRepo) ReadMany(
	ctx context.Context, tx entities.Tx, userID, startID, batchSize int,
) ([]*common.FileReq, error) {
	var files []common.File
	query := "select * from files where user_id = $1 and id > $2 and not is_deleted order by id limit $3"
	if err := tx.SelectContext(ctx, &files, query, userID, startID, batchSize); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	var result = make([]*common.FileReq, 0, len(files))
	for _, row := range files {
		resultRow := common.FileReq{
			ServerID: &row.ID,
			Meta:     row.Meta,
			Value: &common.EncryptionResult{
				Ciphertext: row.EncryptedContent,
				Salt:       row.Salt,
				Nonce:      row.Nonce,
			},
		}
		result = append(result, &resultRow)
	}
	return result, nil
}

func (r *FilesRepo) create(
	ctx context.Context, tx entities.Tx, userID int, data *common.FileReq,
) (int, error) {
	var result common.File
	query := `
		insert into files(user_id, meta, encrypted_content, salt, nonce)
		values ($1, $2, $3, $4, $5)
		returning *
	`
	if err := tx.GetContext(
		ctx, &result, query, userID, data.Meta, data.Value.Ciphertext, data.Value.Salt, data.Value.Nonce,
	); err != nil {
		r.logger.Errorf("Failed to create file: %s", err.Error())
		return 0, err
	}
	r.logger.Infof("File created")
	if err := r.createVersion(ctx, tx, result.ID, data, entities.DefaultVersion); err != nil {
		return 0, err
	}
	return result.ID, nil
}

func (r *FilesRepo) createVersion(
	ctx context.Context, tx entities.Tx, fileID int, data *common.FileReq, version int,
) error {
	query := `
		insert into files_versions(file_id, version, meta, encrypted_content, salt, nonce)
		values ($1, $2, $3, $4, $5, $6)
	`
	if err := tx.ExecContext(
		ctx,
		query,
		fileID,
		version,
		data.Meta,
		data.Value.Ciphertext,
		data.Value.Salt,
		data.Value.Nonce,
	); err != nil {
		r.logger.Errorf("Failed to create file version: %s", err.Error())
		return err
	}
	r.logger.Infof("File version created")
	return nil
}

func (r *FilesRepo) update(ctx context.Context, tx entities.Tx, userID int, data *common.FileReq) error {
	_, version, err := r.Read(ctx, tx, userID, *data.ServerID, true)
	if err != nil {
		return err
	}
	query := `
		update files
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
		r.logger.Errorf("Failed to update file: %s", err.Error())
		return err
	}
	if err := r.createVersion(ctx, tx, *data.ServerID, data, version+1); err != nil {
		return err
	}
	return nil
}
