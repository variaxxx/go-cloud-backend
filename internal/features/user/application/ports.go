package user_app

import (
	user_domain "cloud/internal/features/user/domain"
	"context"
)

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
