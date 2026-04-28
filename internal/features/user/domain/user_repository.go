package user_domain

import "context"

type UserRepository interface {
	Create(
		ctx context.Context,
		username string,
		passwordHash string,
	) (User, error)

	FindByUsername(
		ctx context.Context,
		username string,
	) (User, error)
}
