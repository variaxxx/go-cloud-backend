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
}
