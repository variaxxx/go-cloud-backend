package folder_postgres

import (
	core_errors "cloud/internal/core/errors"
	folder_domain "cloud/internal/features/folder/domain"
	infra_postgres "cloud/internal/infra/postgres"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type folderRepository struct {
	pool infra_postgres.Pool
}

func NewFolderRepository(
	pool infra_postgres.Pool,
) *folderRepository {
	return &folderRepository{
		pool: pool,
	}
}

func (r *folderRepository) Create(
	ctx context.Context,
	name string,
	userID int64,
	parentID *uuid.UUID,
) (folder_domain.Folder, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.GetOperationTimeout())
	defer cancel()

	const query = `
		INSERT INTO cloud.folders (name, user_id, parent_id)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at, name, user_id, parent_id;
	`

	var folder folder_domain.Folder
	var dbParentID *uuid.UUID
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

func (r *folderRepository) Update(
	ctx context.Context,
	id uuid.UUID,
	userID int64,
	name string,
	parentID *uuid.UUID,
) (folder_domain.Folder, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.GetOperationTimeout())
	defer cancel()

	const query = `
		UPDATE cloud.folders
		SET name = $3, parent_id = $4, updated_at = NOW()
		WHERE id = $1 AND user_id = $2
		RETURNING id, created_at, updated_at, name, user_id, parent_id;
	`

	var folder folder_domain.Folder
	var dbParentID *uuid.UUID
	if err := r.pool.QueryRow(ctx, query, id, userID, name, parentID).Scan(
		&folder.ID,
		&folder.CreatedAt,
		&folder.UpdatedAt,
		&folder.Name,
		&folder.UserID,
		&dbParentID,
	); err != nil {
		var pgErr *pgconn.PgError
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return folder_domain.Folder{}, fmt.Errorf("update folder: %w", core_errors.ErrNotFound)
		case errors.As(err, &pgErr) && pgErr.Code == "23505":
			return folder_domain.Folder{}, fmt.Errorf("update folder: %w", core_errors.ErrConflict)
		case errors.As(err, &pgErr) && pgErr.Code == "23503":
			return folder_domain.Folder{}, fmt.Errorf("update folder: %w", core_errors.ErrInvalidArgument)
		default:
			return folder_domain.Folder{}, fmt.Errorf("update folder: %w", err)
		}
	}

	folder.ParentID = dbParentID

	return folder, nil
}

func (r *folderRepository) FindByIDAndUserID(
	ctx context.Context,
	id uuid.UUID,
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
	var dbParentID *uuid.UUID
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

func (r *folderRepository) FindByIDAndUserIDWithPath(
	ctx context.Context,
	id uuid.UUID,
	userID int64,
) (folder_domain.Folder, []folder_domain.Folder, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.GetOperationTimeout())
	defer cancel()

	const query = `
		WITH RECURSIVE ancestry AS (
			SELECT id, created_at, updated_at, name, user_id, parent_id, 0 AS depth
			FROM cloud.folders
			WHERE id = $1 AND user_id = $2

			UNION ALL

			SELECT f.id, f.created_at, f.updated_at, f.name, f.user_id, f.parent_id, a.depth + 1
			FROM cloud.folders f
			JOIN ancestry a ON a.parent_id = f.id
			WHERE f.user_id = $2
		)
		SELECT id, created_at, updated_at, name, user_id, parent_id, depth
		FROM ancestry
		ORDER BY depth DESC;
	`

	rows, err := r.pool.Query(ctx, query, id, userID)
	if err != nil {
		return folder_domain.Folder{}, nil, fmt.Errorf("find folder by id and user id with path: %w", err)
	}
	defer rows.Close()

	var folder folder_domain.Folder
	path := make([]folder_domain.Folder, 0)
	found := false
	for rows.Next() {
		var pathFolder folder_domain.Folder
		var dbParentID *uuid.UUID
		var depth int
		if err := rows.Scan(
			&pathFolder.ID,
			&pathFolder.CreatedAt,
			&pathFolder.UpdatedAt,
			&pathFolder.Name,
			&pathFolder.UserID,
			&dbParentID,
			&depth,
		); err != nil {
			return folder_domain.Folder{}, nil, fmt.Errorf("scan folder path row: %w", err)
		}

		pathFolder.ParentID = dbParentID
		path = append(path, pathFolder)

		if depth == 0 {
			folder = pathFolder
			found = true
		}
	}

	if err := rows.Err(); err != nil {
		return folder_domain.Folder{}, nil, fmt.Errorf("iterate folder path rows: %w", err)
	}

	if !found {
		return folder_domain.Folder{}, nil, fmt.Errorf("find folder by id and user id with path: %w", core_errors.ErrNotFound)
	}

	return folder, path, nil
}

func (r *folderRepository) FindByParentIDAndUserID(
	ctx context.Context,
	parentID *uuid.UUID,
	userID int64,
) ([]folder_domain.Folder, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.GetOperationTimeout())
	defer cancel()

	query := `
		SELECT id, created_at, updated_at, name, user_id, parent_id
		FROM cloud.folders
		WHERE user_id = $1
	`
	args := []any{userID}

	if parentID == nil {
		query += ` AND parent_id IS NULL`
	} else {
		query += ` AND parent_id = $2`
		args = append(args, *parentID)
	}

	query += ` ORDER BY name`

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("find folders by parent id and user id: %w", err)
	}
	defer rows.Close()

	folders := make([]folder_domain.Folder, 0)
	for rows.Next() {
		var folder folder_domain.Folder
		var dbParentID *uuid.UUID
		if err := rows.Scan(
			&folder.ID,
			&folder.CreatedAt,
			&folder.UpdatedAt,
			&folder.Name,
			&folder.UserID,
			&dbParentID,
		); err != nil {
			return nil, fmt.Errorf("scan folder row: %w", err)
		}

		folder.ParentID = dbParentID
		folders = append(folders, folder)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate folder rows: %w", err)
	}

	return folders, nil
}

func (r *folderRepository) WouldCreateCycle(
	ctx context.Context,
	id uuid.UUID,
	parentID uuid.UUID,
	userID int64,
) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.GetOperationTimeout())
	defer cancel()

	const query = `
		WITH RECURSIVE ancestry AS (
			SELECT id, parent_id
			FROM cloud.folders
			WHERE id = $1 AND user_id = $3

			UNION ALL

			SELECT f.id, f.parent_id
			FROM cloud.folders f
			JOIN ancestry a ON a.parent_id = f.id
			WHERE f.user_id = $3
		)
		SELECT EXISTS(
			SELECT 1
			FROM ancestry
			WHERE id = $2
		);
	`

	var wouldCreateCycle bool
	if err := r.pool.QueryRow(ctx, query, parentID, id, userID).Scan(&wouldCreateCycle); err != nil {
		return false, fmt.Errorf("check folder cycle: %w", err)
	}

	return wouldCreateCycle, nil
}

func (r *folderRepository) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.GetOperationTimeout())
	defer cancel()

	const query = `
		DELETE FROM cloud.folders
		WHERE id = $1;
	`

	commandTag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete folder: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("delete folder: %w", core_errors.ErrNotFound)
	}

	return nil
}
