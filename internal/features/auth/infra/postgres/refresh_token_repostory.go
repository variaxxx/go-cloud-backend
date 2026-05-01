package auth_postgres

import (
	core_errors "cloud/internal/core/errors"
	auth_domain "cloud/internal/features/auth/domain"
	infra_postgres "cloud/internal/infra/postgres"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type refreshTokenRepository struct {
	pool infra_postgres.Pool
}

type queryRowScanner interface {
	Scan(dest ...any) error
}

func NewRefreshTokenRepository(
	pool infra_postgres.Pool,
) *refreshTokenRepository {
	return &refreshTokenRepository{
		pool: pool,
	}
}

func (r *refreshTokenRepository) Create(
	ctx context.Context,
	userID int64,
	expiresAt time.Time,
	tokenHash string,
) (auth_domain.RefreshToken, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.GetOperationTimeout())
	defer cancel()

	const query = `
		INSERT INTO cloud.refresh_tokens (expires_at, token_hash, user_id)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, expires_at, revoked_at, token_hash, user_id, replaced_by_id;
	`

	token, err := scanRefreshToken(
		r.pool.QueryRow(ctx, query, expiresAt, tokenHash, userID),
	)
	if err != nil {
		var pgErr *pgconn.PgError
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return auth_domain.RefreshToken{}, fmt.Errorf("create refresh token: %w", core_errors.ErrNotFound)
		case errors.As(err, &pgErr) && pgErr.Code == "23505":
			return auth_domain.RefreshToken{}, fmt.Errorf("create refresh token: %w", core_errors.ErrConflict)
		default:
			return auth_domain.RefreshToken{}, fmt.Errorf("create refresh token: %w", err)
		}
	}

	return token, nil
}

func (r *refreshTokenRepository) Revoke(
	ctx context.Context,
	tokenHash string,
	replacedByID *int64,
) (auth_domain.RefreshToken, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.GetOperationTimeout())
	defer cancel()

	const query = `
		UPDATE cloud.refresh_tokens
		SET revoked_at = NOW(), replaced_by_id = $2
		WHERE token_hash = $1
		RETURNING id, created_at, expires_at, revoked_at, token_hash, user_id, replaced_by_id;
	`

	token, err := scanRefreshToken(
		r.pool.QueryRow(ctx, query, tokenHash, replacedByID),
	)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return auth_domain.RefreshToken{}, fmt.Errorf("revoke refresh token: %w", core_errors.ErrNotFound)
		default:
			return auth_domain.RefreshToken{}, fmt.Errorf("revoke refresh token: %w", err)
		}
	}

	return token, nil
}

func (r *refreshTokenRepository) Rotate(
	ctx context.Context,
	oldTokenHash string,
	newTokenHash string,
	expiresAt time.Time,
	now time.Time,
) (auth_domain.RefreshToken, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.GetOperationTimeout())
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return auth_domain.RefreshToken{}, fmt.Errorf("begin refresh token rotation: %w", err)
	}
	defer tx.Rollback(ctx)

	const selectQuery = `
		SELECT id, created_at, expires_at, revoked_at, token_hash, user_id, replaced_by_id
		FROM cloud.refresh_tokens
		WHERE token_hash = $1
		FOR UPDATE;
	`

	oldToken, err := scanRefreshToken(tx.QueryRow(ctx, selectQuery, oldTokenHash))
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return auth_domain.RefreshToken{}, fmt.Errorf("rotate refresh token: %w", core_errors.ErrNotFound)
		default:
			return auth_domain.RefreshToken{}, fmt.Errorf("rotate refresh token: select current token: %w", err)
		}
	}

	if oldToken.RevokedAt != nil || !oldToken.ExpiresAt.After(now) {
		return auth_domain.RefreshToken{}, fmt.Errorf("rotate refresh token: %w", core_errors.ErrInvalidArgument)
	}

	const insertQuery = `
		INSERT INTO cloud.refresh_tokens (expires_at, token_hash, user_id)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, expires_at, revoked_at, token_hash, user_id, replaced_by_id;
	`

	newToken, err := scanRefreshToken(tx.QueryRow(ctx, insertQuery, expiresAt, newTokenHash, oldToken.UserID))
	if err != nil {
		var pgErr *pgconn.PgError
		switch {
		case errors.As(err, &pgErr) && pgErr.Code == "23505":
			return auth_domain.RefreshToken{}, fmt.Errorf("rotate refresh token: %w", core_errors.ErrConflict)
		default:
			return auth_domain.RefreshToken{}, fmt.Errorf("rotate refresh token: create replacement token: %w", err)
		}
	}

	const updateQuery = `
		UPDATE cloud.refresh_tokens
		SET revoked_at = $2, replaced_by_id = $3
		WHERE token_hash = $1;
	`

	if _, err := tx.Exec(ctx, updateQuery, oldTokenHash, now, newToken.ID); err != nil {
		return auth_domain.RefreshToken{}, fmt.Errorf("rotate refresh token: revoke current token: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return auth_domain.RefreshToken{}, fmt.Errorf("rotate refresh token: commit transaction: %w", err)
	}

	return newToken, nil
}

func scanRefreshToken(row queryRowScanner) (auth_domain.RefreshToken, error) {
	var token auth_domain.RefreshToken
	var revokedAt *time.Time
	var replacedByID *int64

	if err := row.Scan(
		&token.ID,
		&token.CreatedAt,
		&token.ExpiresAt,
		&revokedAt,
		&token.TokenHash,
		&token.UserID,
		&replacedByID,
	); err != nil {
		return auth_domain.RefreshToken{}, err
	}

	token.RevokedAt = revokedAt
	token.ReplacedByID = replacedByID

	return token, nil
}
