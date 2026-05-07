package auth_app

import (
	auth_domain "cloud/internal/features/auth/domain"
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"
)

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

type refreshTokenService struct {
	repository RefreshTokenRepository
	hasher     RefreshTokenHasher
	refreshTTL time.Duration
}

func NewRefreshTokenService(
	repository RefreshTokenRepository,
	hasher RefreshTokenHasher,
	refreshTTL time.Duration,
) *refreshTokenService {
	return &refreshTokenService{
		repository: repository,
		hasher:     hasher,
		refreshTTL: refreshTTL,
	}
}

func (s *refreshTokenService) Issue(
	ctx context.Context,
	userID int64,
) (string, error) {
	expiresAt := time.Now().Add(s.refreshTTL)

	token, err := s.generateToken()
	if err != nil {
		return "", err
	}

	tokenHash := s.hasher.Hash(token)

	_, err = s.repository.Create(
		ctx,
		userID,
		expiresAt,
		tokenHash,
	)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *refreshTokenService) ReplaceOld(
	ctx context.Context,
	oldToken string,
) (string, auth_domain.RefreshToken, error) {
	now := time.Now()
	expiresAt := now.Add(s.refreshTTL)
	newToken, err := s.generateToken()
	if err != nil {
		return "", auth_domain.RefreshToken{}, err
	}

	oldTokenHash := s.hasher.Hash(oldToken)
	newTokenHash := s.hasher.Hash(newToken)

	newTokenInfo, err := s.repository.Rotate(ctx, oldTokenHash, newTokenHash, expiresAt, now)
	if err != nil {
		return "", auth_domain.RefreshToken{}, fmt.Errorf("rotate refresh token: %w", err)
	}

	return newToken, newTokenInfo, nil
}

func (s *refreshTokenService) generateToken() (string, error) {
	b := make([]byte, 32)

	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate refresh token: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}
