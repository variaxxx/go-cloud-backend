package auth_app

import (
	auth_domain "cloud/internal/features/auth/domain"
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"
)

type RefreshTokenService struct {
	repository auth_domain.RefreshTokenRepository
	hasher     auth_domain.RefreshTokenHasher
	refreshTTL time.Duration
}

func NewRefreshTokenService(
	repository auth_domain.RefreshTokenRepository,
	hasher auth_domain.RefreshTokenHasher,
	refreshTTL time.Duration,
) *RefreshTokenService {
	return &RefreshTokenService{
		repository: repository,
		hasher:     hasher,
		refreshTTL: refreshTTL,
	}
}

func (s *RefreshTokenService) Issue(
	ctx context.Context,
	userId int64,
) (string, error) {
	expiresAt := time.Now().Add(s.refreshTTL)

	token, err := s.generateToken()
	if err != nil {
		return "", err
	}

	tokenHash := s.hasher.Hash(token)

	_, err = s.repository.Create(
		ctx,
		userId,
		expiresAt,
		tokenHash,
	)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *RefreshTokenService) ReplaceOld(
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

func (s *RefreshTokenService) generateToken() (string, error) {
	b := make([]byte, 32)

	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate refresh token: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}
