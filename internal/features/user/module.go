package user

import (
	core_http_server "cloud/internal/core/transport/http/server"
	auth_jwt "cloud/internal/features/auth/infra/jwt"
	user_app "cloud/internal/features/user/application"
	user_postgres "cloud/internal/features/user/infra/postgres"
	user_http "cloud/internal/features/user/transport/http"
	infra_postgres "cloud/internal/infra/postgres"
)

type Deps struct {
	Router *core_http_server.APIRouter
	DB     infra_postgres.Pool
}

func Register(
	deps Deps,
) error {
	jwtConfig, err := auth_jwt.NewConfig()
	if err != nil {
		return err
	}

	userRepo := user_postgres.NewUserRepository(deps.DB)
	userService := user_app.NewUserService(userRepo)
	tokenManager := auth_jwt.NewManager(jwtConfig)
	handler := user_http.NewHandler(userService)

	deps.Router.RegisterRoutes(handler.Routes(tokenManager)...)
	return nil
}
