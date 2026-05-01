package auth_app

import (
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
