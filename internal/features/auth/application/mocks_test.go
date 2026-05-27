package auth_app

import (
	auth_domain "cloud/internal/features/auth/domain"
	user_domain "cloud/internal/features/user/domain"
	"context"
	"strconv"
	"time"
)

type mockPasswordHasher struct {
	hashErr          error
	compareErr       error
	comparedHash     string
	comparedPassword string
}

func (f *mockPasswordHasher) Hash(password string) (string, error) {
	if f.hashErr != nil {
		return "", f.hashErr
	}
	return "hashed:" + password, nil
}

func (f *mockPasswordHasher) Compare(hash string, password string) error {
	f.comparedHash = hash
	f.comparedPassword = password
	if f.compareErr != nil {
		return f.compareErr
	}
	return nil
}

type mockRefreshTokenHasher struct{}

func (f *mockRefreshTokenHasher) Hash(token string) string {
	return "hash:" + token
}

type mockTokenManager struct {
	issueErr error
}

func (f *mockTokenManager) Parse(tokenString string) (int64, error) {
	return 0, nil
}

func (f *mockTokenManager) Issue(userID int64) (string, error) {
	if f.issueErr != nil {
		return "", f.issueErr
	}
	return "access:" + strconv.FormatInt(userID, 10), nil
}

type mockRefreshTokenUseCase struct {
	issueErr      error
	replaceErr    error
	replaceUserID int64
	replacedToken string
}

func (f *mockRefreshTokenUseCase) Issue(ctx context.Context, userID int64) (string, error) {
	if f.issueErr != nil {
		return "", f.issueErr
	}
	return "refresh:" + strconv.FormatInt(userID, 10), nil
}

func (f *mockRefreshTokenUseCase) ReplaceOld(ctx context.Context, oldToken string) (string, auth_domain.RefreshToken, error) {
	f.replacedToken = oldToken
	if f.replaceErr != nil {
		return "", auth_domain.RefreshToken{}, f.replaceErr
	}
	return "new-refresh:" + strconv.FormatInt(f.replaceUserID, 10), auth_domain.RefreshToken{UserID: f.replaceUserID}, nil
}

type mockAuthUserRepository struct {
	createErr           error
	findErr             error
	found               user_domain.User
	createdUsername     string
	createdPasswordHash string
	findUsername        string
}

func (f *mockAuthUserRepository) Create(ctx context.Context, username string, passwordHash string) (user_domain.User, error) {
	if f.createErr != nil {
		return user_domain.User{}, f.createErr
	}
	f.createdUsername = username
	f.createdPasswordHash = passwordHash
	return user_domain.User{ID: 100, Username: username, PasswordHash: passwordHash}, nil
}

func (f *mockAuthUserRepository) FindByUsername(ctx context.Context, username string) (user_domain.User, error) {
	if f.findErr != nil {
		return user_domain.User{}, f.findErr
	}
	f.findUsername = username
	return f.found, nil
}

func (f *mockAuthUserRepository) FindByID(ctx context.Context, id int64) (user_domain.User, error) {
	return user_domain.User{ID: id}, nil
}

type mockRefreshTokenRepository struct {
	createErr        error
	rotateErr        error
	rotated          auth_domain.RefreshToken
	createdUserID    int64
	createdExpiresAt time.Time
	createdHash      string
	rotatedOldHash   string
	rotatedNewHash   string
}

func (f *mockRefreshTokenRepository) Create(ctx context.Context, userID int64, expiresAt time.Time, tokenHash string) (auth_domain.RefreshToken, error) {
	if f.createErr != nil {
		return auth_domain.RefreshToken{}, f.createErr
	}
	f.createdUserID = userID
	f.createdExpiresAt = expiresAt
	f.createdHash = tokenHash
	return auth_domain.RefreshToken{UserID: userID, ExpiresAt: expiresAt, TokenHash: tokenHash}, nil
}

func (f *mockRefreshTokenRepository) Revoke(ctx context.Context, tokenHash string, replacedByID *int64) (auth_domain.RefreshToken, error) {
	return auth_domain.RefreshToken{}, nil
}

func (f *mockRefreshTokenRepository) Rotate(ctx context.Context, oldTokenHash string, newTokenHash string, expiresAt time.Time, now time.Time) (auth_domain.RefreshToken, error) {
	if f.rotateErr != nil {
		return auth_domain.RefreshToken{}, f.rotateErr
	}
	f.rotatedOldHash = oldTokenHash
	f.rotatedNewHash = newTokenHash
	return f.rotated, nil
}
