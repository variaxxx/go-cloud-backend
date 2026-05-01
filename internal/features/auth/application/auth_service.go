package auth_app

import (
	core_errors "cloud/internal/core/errors"
	auth_domain "cloud/internal/features/auth/domain"
	user_domain "cloud/internal/features/user/domain"
	"context"
	"fmt"
	"strings"
)

type authService struct {
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
) *authService {
	return &authService{
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

func (s *authService) Register(
	ctx context.Context,
	username string,
	password string,
) (Tokens, error) {
	username = strings.TrimSpace(username)
	if username == "" {
		return Tokens{}, fmt.Errorf("register user: %w", core_errors.ErrInvalidArgument)
	}

	passwordHash, err := s.hasher.Hash(password)
	if err != nil {
		return Tokens{}, fmt.Errorf("hash password: %w", err)
	}

	user, err := s.userRepository.Create(ctx, username, passwordHash)
	if err != nil {
		return Tokens{}, fmt.Errorf("register user: %w", err)
	}

	tokens, err := s.issueTokens(ctx, user.ID)
	if err != nil {
		return Tokens{}, err
	}

	return tokens, nil
}

func (s *authService) Login(
	ctx context.Context,
	username string,
	password string,
) (Tokens, error) {
	username = strings.TrimSpace(username)
	if username == "" {
		return Tokens{}, fmt.Errorf("find user by username: %w", core_errors.ErrInvalidArgument)
	}

	user, err := s.userRepository.FindByUsername(ctx, username)
	if err != nil {
		return Tokens{}, fmt.Errorf("find user by username: %w", err)
	}

	if err := s.hasher.Compare(user.PasswordHash, password); err != nil {
		return Tokens{}, fmt.Errorf("%w: invalid password: %w", core_errors.ErrInvalidArgument, err)
	}

	tokens, err := s.issueTokens(ctx, user.ID)
	if err != nil {
		return Tokens{}, err
	}

	return tokens, nil
}

func (s *authService) RefreshTokens(
	ctx context.Context,
	refreshToken string,
) (Tokens, error) {
	newRefresh, newRefreshInfo, err := s.refreshTokenService.ReplaceOld(ctx, refreshToken)
	if err != nil {
		return Tokens{}, fmt.Errorf("refresh tokens: %w", err)
	}

	accessToken, err := s.tokenManager.Issue(newRefreshInfo.UserID)
	if err != nil {
		return Tokens{}, fmt.Errorf("access token issue: %w", err)
	}

	return Tokens{
		AccessToken:  accessToken,
		RefreshToken: newRefresh,
	}, nil
}

func (s *authService) issueTokens(
	ctx context.Context,
	userID int64,
) (Tokens, error) {
	accessToken, err := s.tokenManager.Issue(userID)
	if err != nil {
		return Tokens{}, fmt.Errorf("access token issue: %w", err)
	}
	refreshToken, err := s.refreshTokenService.Issue(ctx, userID)
	if err != nil {
		return Tokens{}, fmt.Errorf("refresh token issue: %w", err)
	}

	return Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
