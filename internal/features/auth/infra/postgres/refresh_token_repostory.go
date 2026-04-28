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

type RefreshTokenRepository struct {
	pool infra_postgres.Pool
}

func NewRefreshTokenRepository(
	pool infra_postgres.Pool,
) *RefreshTokenRepository {
	return &RefreshTokenRepository{
		pool: pool,
	}
}

func (r *RefreshTokenRepository) Create(
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

	var token auth_domain.RefreshToken
	var revokedAt *time.Time
	var replacedByID *int64
	if err := r.pool.QueryRow(ctx, query, expiresAt, tokenHash, userID).Scan(
		&token.ID,
		&token.CreatedAt,
		&token.ExpiresAt,
		&revokedAt,
		&token.TokenHash,
		&token.UserID,
		&replacedByID,
	); err != nil {
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

	token.RevokedAt = revokedAt
	token.ReplacedByID = replacedByID

	return token, nil
}
