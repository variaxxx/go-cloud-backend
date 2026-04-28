package auth_app

import (
	auth_domain "cloud/internal/features/auth/domain"
	"context"
)

type AuthUseCase interface {
	Register(
		ctx context.Context,
		username string,
		password string,
	) (Tokens, error)

	Login(
		ctx context.Context,
		username string,
		password string,
	) (Tokens, error)

	RefreshTokens(
		ctx context.Context,
		token string,
	) (Tokens, error)
}

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
