package auth

import (
	core_http_server "cloud/internal/core/transport/http/server"
	auth_app "cloud/internal/features/auth/application"
	auth_jwt "cloud/internal/features/auth/infra/jwt"
	auth_password "cloud/internal/features/auth/infra/password"
	auth_postgres "cloud/internal/features/auth/infra/postgres"
	auth_refresh_token "cloud/internal/features/auth/infra/refresh_token"
	auth_http "cloud/internal/features/auth/transport/http"
	user_postgres "cloud/internal/features/user/infra/postgres"
	infra_postgres "cloud/internal/infra/postgres"
)

type Deps struct {
	Router *core_http_server.APIRouter
	DB     infra_postgres.Pool
}

func Register(
	deps Deps,
) error {
	config, err := NewConfig()
	if err != nil {
		return err
	}

	jwtConfig, err := auth_jwt.NewConfig()
	if err != nil {
		return err
	}

	userRepo := user_postgres.NewUserRepository(deps.DB)
	refreshTokenRepo := auth_postgres.NewRefreshTokenRepository(deps.DB)
	tokenManager := auth_jwt.NewManager(jwtConfig)
	hasher := auth_password.NewHasher()
	refreshTokenHasher := auth_refresh_token.NewHasher()

	refreshTokenService := auth_app.NewRefreshTokenService(refreshTokenRepo, refreshTokenHasher, config.RefreshTtl)
	authService := auth_app.NewAuthService(userRepo, tokenManager, hasher, refreshTokenService)

	handler := auth_http.NewHandler(authService, config.RefreshTtl)

	deps.Router.RegisterRoutes(handler.Routes()...)
	return nil
}
