package repositories

import (
	"context"
	"database/sql"
	"errors"

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
	if data.ServerID == nil {
		r.logger.Infof("Creating new note for user: %d", userID)
		return r.create(ctx, tx, userID, data)
	}
	r.logger.Infof("Updating note %d for user: %d", data.ServerID, userID)
	return *data.ServerID, r.update(ctx, tx, userID, data)
}

func (r *NotesRepo) Read(
	ctx context.Context, tx entities.Tx, userID int, rowID int, lock bool,
) (*common.NoteReq, int, error) {
	var note common.Note
	query := "select * from notes where id = $1 and user_id = $2 and not is_deleted"
	if lock {
		query = query + " for update"
	}
	if err := tx.GetContext(ctx, &note, query, rowID, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, 0, entities.ErrDoesntExist
		}
		return nil, 0, err
	}
	var result common.NoteReq
	result.ServerID = &rowID
	result.Meta = note.Meta
	result.Value = &common.EncryptionResult{
		Ciphertext: note.EncryptedContent,
		Salt:       note.Salt,
		Nonce:      note.Nonce,
	}
	return &result, note.Version, nil
}

func (r *NotesRepo) Delete(ctx context.Context, tx entities.Tx, userID int, rowID int) error {
	_, _, err := r.Read(ctx, tx, userID, rowID, true)
	if err != nil {
		return err
	}
	query := "update notes set is_deleted = true where id = $1"
	if err := tx.ExecContext(ctx, query, rowID); err != nil {
		r.logger.Errorf("Failed to delete note: %s", err.Error())
		return err
	}
	return nil
}

func (r *NotesRepo) ReadMany(
	ctx context.Context, tx entities.Tx, userID, startID, batchSize int,
) ([]*common.NoteReq, error) {
	var notes []common.Note
	query := "select * from notes where user_id = $1 and id > $2 and not is_deleted order by id limit $3"
	if err := tx.SelectContext(ctx, &notes, query, userID, startID, batchSize); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	var result = make([]*common.NoteReq, 0, len(notes))
	for _, row := range notes {
		resultRow := common.NoteReq{
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
	if err := r.createVersion(ctx, tx, result.ID, data, entities.DefaultVersion); err != nil {
		return 0, err
	}
	return result.ID, nil
}

func (r *NotesRepo) createVersion(
	ctx context.Context, tx entities.Tx, noteID int, data *common.NoteReq, version int,
) error {
	query := `
		insert into notes_versions(note_id, version, meta, encrypted_content, salt, nonce)
		values ($1, $2, $3, $4, $5, $6)
	`
	if err := tx.ExecContext(
		ctx,
		query,
		noteID,
		version,
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

func (r *NotesRepo) update(ctx context.Context, tx entities.Tx, userID int, data *common.NoteReq) error {
	_, version, err := r.Read(ctx, tx, userID, *data.ServerID, true)
	if err != nil {
		return err
	}
	query := `
		update notes
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
		r.logger.Errorf("Failed to update note: %s", err.Error())
		return err
	}
	if err := r.createVersion(ctx, tx, *data.ServerID, data, version+1); err != nil {
		return err
	}
	return nil
}
