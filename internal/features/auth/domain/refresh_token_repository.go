package auth_domain

import (
	"context"
	"time"
)

type RefreshTokenRepository interface {
	Create(
		ctx context.Context,
		userID int64,
		expiresAt time.Time,
		tokenHash string,
	) (RefreshToken, error)

	Revoke(
		ctx context.Context,
		tokenHash string,
		replacedByID *int64,
	) (RefreshToken, error)

	FindByHash(
		ctx context.Context,
		tokenHash string,
	) (RefreshToken, error)
}
