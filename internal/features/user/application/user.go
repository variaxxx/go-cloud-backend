package user_app

import (
	user_domain "cloud/internal/features/user/domain"
	"context"
	"fmt"
)

type UserUseCase interface {
	GetMe(
		ctx context.Context,
		userID int64,
	) (user_domain.User, error)
}

type userService struct {
	userRepo UserRepository
}

func NewUserService(
	userRepo UserRepository,
) *userService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetMe(
	ctx context.Context,
	userID int64,
) (user_domain.User, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return user_domain.User{}, fmt.Errorf("find user by id: %w", err)
	}

	return user, nil
}
