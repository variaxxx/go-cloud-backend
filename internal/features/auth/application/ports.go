package auth_app

import (
	auth_domain "cloud/internal/features/auth/domain"
	user_domain "cloud/internal/features/user/domain"
	"context"
	"time"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash string, password string) error
}

type RefreshTokenHasher interface {
	Hash(token string) string
}

type TokenManager interface {
	Parse(tokenString string) (int64, error)
	Issue(userID int64) (string, error)
}

type RefreshTokenRepository interface {
	Create(
		ctx context.Context,
		userID int64,
		expiresAt time.Time,
		tokenHash string,
	) (auth_domain.RefreshToken, error)

	Revoke(
		ctx context.Context,
		tokenHash string,
		replacedByID *int64,
	) (auth_domain.RefreshToken, error)

	Rotate(
		ctx context.Context,
		oldTokenHash string,
		newTokenHash string,
		expiresAt time.Time,
		now time.Time,
	) (auth_domain.RefreshToken, error)
}

type UserRepository interface {
	Create(
		ctx context.Context,
		username string,
		passwordHash string,
	) (user_domain.User, error)

	FindByUsername(
		ctx context.Context,
		username string,
	) (user_domain.User, error)

	FindByID(
		ctx context.Context,
		id int64,
	) (user_domain.User, error)
}
