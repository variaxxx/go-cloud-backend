package user_postgres

import (
	core_errors "cloud/internal/core/errors"
	user_domain "cloud/internal/features/user/domain"
	infra_postgres "cloud/internal/infra/postgres"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type userRepository struct {
	pool infra_postgres.Pool
}

func NewUserRepository(
	pool infra_postgres.Pool,
) *userRepository {
	return &userRepository{
		pool: pool,
	}
}

func (r *userRepository) Create(
	ctx context.Context,
	username string,
	passwordHash string,
) (user_domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.GetOperationTimeout())
	defer cancel()

	const query = `
		INSERT INTO cloud.users (username, password_hash)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at, username, password_hash
	`

	var user user_domain.User
	if err := r.pool.QueryRow(ctx, query, username, passwordHash).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Username,
		&user.PasswordHash,
	); err != nil {
		var pgErr *pgconn.PgError
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return user_domain.User{}, fmt.Errorf("create user: %w", core_errors.ErrNotFound)
		case errors.As(err, &pgErr) && pgErr.Code == "23505":
			return user_domain.User{}, fmt.Errorf("create user: %w", core_errors.ErrConflict)
		default:
			return user_domain.User{}, fmt.Errorf("create user: %w", err)
		}
	}

	return user, nil
}

func (r *userRepository) FindByUsername(
	ctx context.Context,
	username string,
) (user_domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.GetOperationTimeout())
	defer cancel()

	const query = `
		SELECT id, created_at, updated_at, username, password_hash
		FROM cloud.users
		WHERE username = $1;
	`

	var user user_domain.User
	if err := r.pool.QueryRow(ctx, query, username).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Username,
		&user.PasswordHash,
	); err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return user_domain.User{}, fmt.Errorf("find user by username: %w", core_errors.ErrNotFound)
		default:
			return user_domain.User{}, fmt.Errorf("find user by username: %w", err)
		}
	}

	return user, nil
}

func (r *userRepository) FindByID(
	ctx context.Context,
	id int64,
) (user_domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.GetOperationTimeout())
	defer cancel()

	const query = `
		SELECT id, created_at, updated_at, username, password_hash
		FROM cloud.users
		WHERE id = $1;
	`

	var user user_domain.User
	if err := r.pool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Username,
		&user.PasswordHash,
	); err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return user_domain.User{}, fmt.Errorf("find user by id: %w", core_errors.ErrNotFound)
		default:
			return user_domain.User{}, fmt.Errorf("find user by id: %w", err)
		}
	}

	return user, nil
}
