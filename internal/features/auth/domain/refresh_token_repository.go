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

	Rotate(
		ctx context.Context,
		oldTokenHash string,
		newTokenHash string,
		expiresAt time.Time,
		now time.Time,
	) (RefreshToken, error)
}
