package auth_app

import (
	core_errors "cloud/internal/core/errors"
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
	oldTokenHash := s.hasher.Hash(oldToken)
	oldTokenInfo, err := s.repository.FindByHash(ctx, oldTokenHash)
	if err != nil {
		return "", auth_domain.RefreshToken{}, fmt.Errorf("get old token info: %w", err)
	}

	now := time.Now()
	if oldTokenInfo.RevokedAt != nil {
		return "", auth_domain.RefreshToken{}, fmt.Errorf("%w: refresh token has been revoked", core_errors.ErrInvalidArgument)
	}

	if !oldTokenInfo.ExpiresAt.After(now) {
		return "", auth_domain.RefreshToken{}, fmt.Errorf("%w: refresh token has expired", core_errors.ErrInvalidArgument)
	}

	expiresAt := now.Add(s.refreshTTL)
	newToken, err := s.generateToken()
	if err != nil {
		return "", auth_domain.RefreshToken{}, err
	}
	newTokenHash := s.hasher.Hash(newToken)

	newTokenInfo, err := s.repository.Create(
		ctx,
		oldTokenInfo.UserID,
		expiresAt,
		newTokenHash,
	)
	if err != nil {
		return "", auth_domain.RefreshToken{}, err
	}

	_, err = s.repository.Revoke(ctx, oldTokenHash, &newTokenInfo.ID)
	if err != nil {
		return "", auth_domain.RefreshToken{}, fmt.Errorf("revoke old token: %w", err)
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
