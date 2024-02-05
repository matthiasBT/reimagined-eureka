package repositories

import (
	"context"

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
	r.logger.Infof("Creating new file for user: %d", userID)
	return r.create(ctx, tx, userID, data)
}

func (r *FilesRepo) Read(ctx context.Context, tx entities.Tx, userID int, rowId int) (*common.FileReq, error) {
	panic("implement me!")
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
	if err := r.createVersion(ctx, tx, result.ID, data); err != nil {
		return 0, err
	}
	return result.ID, nil
}

func (r *FilesRepo) createVersion(
	ctx context.Context, tx entities.Tx, fileID int, data *common.FileReq,
) error {
	query := `
		insert into files_versions(file_id, version, meta, encrypted_content, salt, nonce)
		values ($1, $2, $3, $4, $5, $6)
	`
	if err := tx.ExecContext(
		ctx,
		query,
		fileID,
		entities.DefaultVersion,
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

func (r *FilesRepo) update(ctx context.Context, tx entities.Tx, userID int, data *common.FileReq) (int, error) {
	panic("Implement me!")
}
