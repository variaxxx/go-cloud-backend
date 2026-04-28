package auth_app

import (
	auth_domain "cloud/internal/features/auth/domain"
	user_domain "cloud/internal/features/user/domain"
	"context"
	"fmt"
)

type AuthService struct {
	userRepository user_domain.UserRepository
	tokenManager   auth_domain.TokenManager
	hasher         auth_domain.PasswordHasher

	refreshTokenService RefreshTokenUseCase
}

func NewAuthService(
	userRepository user_domain.UserRepository,
	tokenManager auth_domain.TokenManager,
	hasher auth_domain.PasswordHasher,
	refreshTokenService RefreshTokenUseCase,
) *AuthService {
	return &AuthService{
		userRepository:      userRepository,
		tokenManager:        tokenManager,
		hasher:              hasher,
		refreshTokenService: refreshTokenService,
	}
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

func (s *AuthService) Register(
	ctx context.Context,
	username string,
	password string,
) (Tokens, error) {
	passwordHash, err := s.hasher.Hash(password)
	if err != nil {
		return Tokens{}, fmt.Errorf("hash password: %w", err)
	}

	user, err := s.userRepository.Create(ctx, username, passwordHash)
	if err != nil {
		return Tokens{}, fmt.Errorf("register user: %w", err)
	}

	accessToken, err := s.tokenManager.Issue(user.ID)
	if err != nil {
		return Tokens{}, fmt.Errorf("access token issue: %w", err)
	}
	refreshToken, err := s.refreshTokenService.Issue(ctx, user.ID)
	if err != nil {
		return Tokens{}, fmt.Errorf("refresh token issue: %w", err)
	}

	tokens := Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return tokens, nil
}
