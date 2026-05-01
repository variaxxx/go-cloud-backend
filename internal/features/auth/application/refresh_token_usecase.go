package auth_app

import (
	auth_domain "cloud/internal/features/auth/domain"
	"context"
)

type RefreshTokenUseCase interface {
	Issue(
		ctx context.Context,
		userID int64,
	) (string, error)

	ReplaceOld(
		ctx context.Context,
		oldToken string,
	) (string, auth_domain.RefreshToken, error)
}
