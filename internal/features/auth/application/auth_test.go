package auth_app

import (
	core_errors "cloud/internal/core/errors"
	auth_domain "cloud/internal/features/auth/domain"
	user_domain "cloud/internal/features/user/domain"
	"context"
	"errors"
	"testing"
	"time"
)

func TestAuthServiceRegister(t *testing.T) {
	ctx := context.Background()
	userRepo := &mockAuthUserRepository{}
	tokenManager := &mockTokenManager{}
	refreshTokens := &mockRefreshTokenUseCase{}
	service := NewAuthService(userRepo, tokenManager, &mockPasswordHasher{}, refreshTokens)

	got, err := service.Register(ctx, "  alice  ", "secret")
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	if userRepo.createdUsername != "alice" {
		t.Fatalf("created username = %q, want alice", userRepo.createdUsername)
	}
	if userRepo.createdPasswordHash != "hashed:secret" {
		t.Fatalf("created password hash = %q, want hashed:secret", userRepo.createdPasswordHash)
	}
	if got.AccessToken != "access:100" || got.RefreshToken != "refresh:100" {
		t.Fatalf("tokens = %+v, want access:100/refresh:100", got)
	}
}

func TestAuthServiceRegisterRejectsBlankUsername(t *testing.T) {
	service := NewAuthService(&mockAuthUserRepository{}, &mockTokenManager{}, &mockPasswordHasher{}, &mockRefreshTokenUseCase{})

	_, err := service.Register(context.Background(), "   ", "secret")
	if !errors.Is(err, core_errors.ErrInvalidArgument) {
		t.Fatalf("Register() error = %v, want ErrInvalidArgument", err)
	}
}

func TestAuthServiceLogin(t *testing.T) {
	ctx := context.Background()

	t.Run("issues tokens after password check", func(t *testing.T) {
		hasher := &mockPasswordHasher{}
		userRepo := &mockAuthUserRepository{found: user_domain.User{ID: 25, Username: "alice", PasswordHash: "hashed:secret"}}
		service := NewAuthService(userRepo, &mockTokenManager{}, hasher, &mockRefreshTokenUseCase{})

		got, err := service.Login(ctx, " alice ", "secret")
		if err != nil {
			t.Fatalf("Login() error = %v", err)
		}

		if userRepo.findUsername != "alice" {
			t.Fatalf("find username = %q, want alice", userRepo.findUsername)
		}
		if hasher.comparedHash != "hashed:secret" || hasher.comparedPassword != "secret" {
			t.Fatalf("compare args = (%q, %q), want hashed:secret/secret", hasher.comparedHash, hasher.comparedPassword)
		}
		if got.AccessToken != "access:25" || got.RefreshToken != "refresh:25" {
			t.Fatalf("tokens = %+v, want access:25/refresh:25", got)
		}
	})

	t.Run("rejects invalid password", func(t *testing.T) {
		hasher := &mockPasswordHasher{compareErr: errors.New("mismatch")}
		userRepo := &mockAuthUserRepository{found: user_domain.User{ID: 25, Username: "alice", PasswordHash: "hashed:secret"}}
		service := NewAuthService(userRepo, &mockTokenManager{}, hasher, &mockRefreshTokenUseCase{})

		_, err := service.Login(ctx, "alice", "bad")
		if !errors.Is(err, core_errors.ErrInvalidArgument) {
			t.Fatalf("Login() error = %v, want ErrInvalidArgument", err)
		}
	})
}

func TestAuthServiceRefreshTokens(t *testing.T) {
	ctx := context.Background()
	refreshTokens := &mockRefreshTokenUseCase{replaceUserID: 77}
	service := NewAuthService(&mockAuthUserRepository{}, &mockTokenManager{}, &mockPasswordHasher{}, refreshTokens)

	got, err := service.RefreshTokens(ctx, "old-refresh")
	if err != nil {
		t.Fatalf("RefreshTokens() error = %v", err)
	}

	if refreshTokens.replacedToken != "old-refresh" {
		t.Fatalf("replaced token = %q, want old-refresh", refreshTokens.replacedToken)
	}
	if got.AccessToken != "access:77" || got.RefreshToken != "new-refresh:77" {
		t.Fatalf("tokens = %+v, want access:77/new-refresh:77", got)
	}
}

func TestRefreshTokenService(t *testing.T) {
	ctx := context.Background()
	ttl := 2 * time.Hour

	t.Run("issue hashes token and persists expiration", func(t *testing.T) {
		repo := &mockRefreshTokenRepository{}
		service := NewRefreshTokenService(repo, &mockRefreshTokenHasher{}, ttl)

		token, err := service.Issue(ctx, 11)
		if err != nil {
			t.Fatalf("Issue() error = %v", err)
		}

		if token == "" {
			t.Fatal("Issue() token is empty")
		}
		if repo.createdUserID != 11 {
			t.Fatalf("created userID = %d, want 11", repo.createdUserID)
		}
		if repo.createdHash != "hash:"+token {
			t.Fatalf("created hash = %q, want hash of token", repo.createdHash)
		}
		if time.Until(repo.createdExpiresAt) <= time.Hour {
			t.Fatalf("created expiration = %s, want roughly %s in future", repo.createdExpiresAt, ttl)
		}
	})

	t.Run("replace old rotates hashed tokens", func(t *testing.T) {
		repo := &mockRefreshTokenRepository{rotated: auth_domain.RefreshToken{ID: 5, UserID: 12}}
		service := NewRefreshTokenService(repo, &mockRefreshTokenHasher{}, ttl)

		newToken, info, err := service.ReplaceOld(ctx, "old-token")
		if err != nil {
			t.Fatalf("ReplaceOld() error = %v", err)
		}

		if newToken == "" {
			t.Fatal("ReplaceOld() new token is empty")
		}
		if repo.rotatedOldHash != "hash:old-token" {
			t.Fatalf("old hash = %q, want hash:old-token", repo.rotatedOldHash)
		}
		if repo.rotatedNewHash != "hash:"+newToken {
			t.Fatalf("new hash = %q, want hash of new token", repo.rotatedNewHash)
		}
		if info.UserID != 12 {
			t.Fatalf("rotated token userID = %d, want 12", info.UserID)
		}
	})
}
