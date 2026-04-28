package auth_app

import "context"

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
}

type RefreshTokenUseCase interface {
	Issue(
		ctx context.Context,
		userID int64,
	) (string, error)
}
