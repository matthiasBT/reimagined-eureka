package repositories

import (
	"context"

	"reimagined_eureka/internal/common"
	"reimagined_eureka/internal/server/entities"
	"reimagined_eureka/internal/server/infra/logging"
)

type NotesRepo struct {
	logger  logging.ILogger
	storage entities.Storage
}

func NewNotesRepo(logger logging.ILogger, storage entities.Storage) *NotesRepo {
	return &NotesRepo{
		logger:  logger,
		storage: storage,
	}
}

func (r *NotesRepo) Write(ctx context.Context, tx entities.Tx, userID int, data *common.NoteReq) (int, error) {
	r.logger.Infof("Creating new note for user: %d", userID)
	return r.create(ctx, tx, userID, data)
}

func (r *NotesRepo) Read(ctx context.Context, tx entities.Tx, userID int, rowId int) (*common.NoteReq, error) {
	panic("implement me!")
}

func (r *NotesRepo) create(
	ctx context.Context, tx entities.Tx, userID int, data *common.NoteReq,
) (int, error) {
	var result common.Note
	query := `
		insert into notes(user_id, meta, encrypted_content, salt, nonce)
		values ($1, $2, $3, $4, $5)
		returning *
	`
	if err := tx.GetContext(
		ctx, &result, query, userID, data.Meta, data.Value.Ciphertext, data.Value.Salt, data.Value.Nonce,
	); err != nil {
		r.logger.Errorf("Failed to create note: %s", err.Error())
		return 0, err
	}
	r.logger.Infof("Note created")
	if err := r.createVersion(ctx, tx, result.ID, data); err != nil {
		return 0, err
	}
	return result.ID, nil
}

func (r *NotesRepo) createVersion(
	ctx context.Context, tx entities.Tx, noteID int, data *common.NoteReq,
) error {
	query := `
		insert into notes_versions(note_id, version, meta, encrypted_content, salt, nonce)
		values ($1, $2, $3, $4, $5, $6)
	`
	if err := tx.ExecContext(
		ctx,
		query,
		noteID,
		entities.DefaultVersion,
		data.Meta,
		data.Value.Ciphertext,
		data.Value.Salt,
		data.Value.Nonce,
	); err != nil {
		r.logger.Errorf("Failed to create note version: %s", err.Error())
		return err
	}
	r.logger.Infof("Note version created")
	return nil
}

func (r *NotesRepo) update(ctx context.Context, tx entities.Tx, userID int, data *common.NoteReq) (int, error) {
	panic("Implement me!")
}
