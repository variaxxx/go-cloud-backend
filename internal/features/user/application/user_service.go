package user_app

import (
	user_domain "cloud/internal/features/user/domain"
	"context"
	"fmt"
)

type UserService struct {
	userRepo user_domain.UserRepository
}

func NewUserService(
	userRepo user_domain.UserRepository,
) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) GetMe(
	ctx context.Context,
	userID int64,
) (user_domain.User, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return user_domain.User{}, fmt.Errorf("find user by id: %w", err)
	}

	return user, nil
}
