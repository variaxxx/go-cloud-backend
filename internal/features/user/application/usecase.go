package user_app

import (
	user_domain "cloud/internal/features/user/domain"
	"context"
)

type UserUseCase interface {
	GetMe(
		ctx context.Context,
		userID int64,
	) (user_domain.User, error)
}
