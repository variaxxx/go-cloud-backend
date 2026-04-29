package folder_postgres

import (
	core_errors "cloud/internal/core/errors"
	folder_domain "cloud/internal/features/folder/domain"
	infra_postgres "cloud/internal/infra/postgres"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type FolderRepository struct {
	pool infra_postgres.Pool
}

func NewFolderRepository(
	pool infra_postgres.Pool,
) *FolderRepository {
	return &FolderRepository{
		pool: pool,
	}
}

func (r *FolderRepository) Create(
	ctx context.Context,
	name string,
	userID int64,
	parentID *int64,
) (folder_domain.Folder, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.GetOperationTimeout())
	defer cancel()

	const query = `
		INSERT INTO cloud.folders (name, user_id, parent_id)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at, name, user_id, parent_id;
	`

	var folder folder_domain.Folder
	var dbParentID *int64
	if err := r.pool.QueryRow(ctx, query, name, userID, parentID).Scan(
		&folder.ID,
		&folder.CreatedAt,
		&folder.UpdatedAt,
		&folder.Name,
		&folder.UserID,
		&dbParentID,
	); err != nil {
		var pgErr *pgconn.PgError
		switch {
		case errors.As(err, &pgErr) && pgErr.Code == "23505":
			return folder_domain.Folder{}, fmt.Errorf("create folder: %w", core_errors.ErrConflict)
		case errors.As(err, &pgErr) && pgErr.Code == "23503":
			return folder_domain.Folder{}, fmt.Errorf("create folder: %w", core_errors.ErrInvalidArgument)
		default:
			return folder_domain.Folder{}, fmt.Errorf("create folder: %w", err)
		}
	}

	folder.ParentID = dbParentID

	return folder, nil
}

func (r *FolderRepository) FindByIDAndUserID(
	ctx context.Context,
	id int64,
	userID int64,
) (folder_domain.Folder, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.GetOperationTimeout())
	defer cancel()

	const query = `
		SELECT id, created_at, updated_at, name, user_id, parent_id
		FROM cloud.folders
		WHERE id = $1 AND user_id = $2;
	`

	var folder folder_domain.Folder
	var dbParentID *int64
	if err := r.pool.QueryRow(ctx, query, id, userID).Scan(
		&folder.ID,
		&folder.CreatedAt,
		&folder.UpdatedAt,
		&folder.Name,
		&folder.UserID,
		&dbParentID,
	); err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return folder_domain.Folder{}, fmt.Errorf("find folder by id and user id: %w", core_errors.ErrNotFound)
		default:
			return folder_domain.Folder{}, fmt.Errorf("find folder by id and user id: %w", err)
		}
	}

	folder.ParentID = dbParentID

	return folder, nil
}
