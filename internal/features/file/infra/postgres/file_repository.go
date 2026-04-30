package file_postgres

import (
	core_errors "cloud/internal/core/errors"
	file_domain "cloud/internal/features/file/domain"
	infra_postgres "cloud/internal/infra/postgres"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

type FileRepository struct {
	pool infra_postgres.Pool
}

func NewFileRepository(
	pool infra_postgres.Pool,
) *FileRepository {
	return &FileRepository{
		pool: pool,
	}
}

func (r *FileRepository) Create(
	ctx context.Context,
	filename string,
	mimetype *string,
	status file_domain.FileStatus,
	storagePath string,
	sizeBytes int64,
	userID int64,
	folderID *uuid.UUID,
) (file_domain.File, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.GetOperationTimeout())
	defer cancel()

	const query = `
		INSERT INTO cloud.files (filename, mimetype, status, storage_path, size_bytes, user_id, folder_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at, deleted_at, filename, mimetype, status, storage_path, size_bytes, user_id, folder_id;
	`

	var file file_domain.File
	var deletedAt *time.Time
	var dbMimetype *string
	var dbFolderID *uuid.UUID
	if err := r.pool.QueryRow(ctx, query, filename, mimetype, status, storagePath, sizeBytes, userID, folderID).Scan(
		&file.ID,
		&file.CreatedAt,
		&file.UpdatedAt,
		&deletedAt,
		&file.Filename,
		&dbMimetype,
		&file.Status,
		&file.StoragePath,
		&file.SizeBytes,
		&file.UserID,
		&dbFolderID,
	); err != nil {
		var pgErr *pgconn.PgError
		switch {
		case errors.As(err, &pgErr) && pgErr.Code == "23505":
			return file_domain.File{}, fmt.Errorf("create file: %w", core_errors.ErrConflict)
		case errors.As(err, &pgErr) && pgErr.Code == "23503":
			return file_domain.File{}, fmt.Errorf("create file: %w", core_errors.ErrInvalidArgument)
		default:
			return file_domain.File{}, fmt.Errorf("create file: %w", err)
		}
	}

	file.DeletedAt = deletedAt
	file.Mimetype = dbMimetype
	file.FolderID = dbFolderID

	return file, nil
}

func (r *FileRepository) FindByFolderIDAndUserID(
	ctx context.Context,
	folderID *uuid.UUID,
	userID int64,
) ([]file_domain.File, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.GetOperationTimeout())
	defer cancel()

	query := `
		SELECT id, created_at, updated_at, deleted_at, filename, mimetype, status, storage_path, size_bytes, user_id, folder_id
		FROM cloud.files
		WHERE user_id = $1 AND deleted_at IS NULL
	`
	args := []any{userID}

	if folderID == nil {
		query += ` AND folder_id IS NULL`
	} else {
		query += ` AND folder_id = $2`
		args = append(args, *folderID)
	}

	query += ` ORDER BY filename`

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("find files by folder id and user id: %w", err)
	}
	defer rows.Close()

	files := make([]file_domain.File, 0)
	for rows.Next() {
		var file file_domain.File
		var deletedAt *time.Time
		var dbMimetype *string
		var dbFolderID *uuid.UUID
		if err := rows.Scan(
			&file.ID,
			&file.CreatedAt,
			&file.UpdatedAt,
			&deletedAt,
			&file.Filename,
			&dbMimetype,
			&file.Status,
			&file.StoragePath,
			&file.SizeBytes,
			&file.UserID,
			&dbFolderID,
		); err != nil {
			return nil, fmt.Errorf("scan file row: %w", err)
		}

		file.DeletedAt = deletedAt
		file.Mimetype = dbMimetype
		file.FolderID = dbFolderID
		files = append(files, file)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate file rows: %w", err)
	}

	return files, nil
}
