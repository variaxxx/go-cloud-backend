package user_http

import (
	core_http_auth "cloud/internal/core/transport/http/auth"
	core_http_middleware "cloud/internal/core/transport/http/middleware"
	core_http_server "cloud/internal/core/transport/http/server"
	user_app "cloud/internal/features/user/application"
	"net/http"
)

type handler struct {
	userService user_app.UserUseCase
}

func NewHandler(
	userService user_app.UserUseCase,
) *handler {
	return &handler{
		userService: userService,
	}
}

func (h *handler) Routes(
	tokenParser core_http_auth.TokenParser,
) []core_http_server.Route {
	return []core_http_server.Route{
		core_http_server.NewRoute(
			http.MethodGet,
			"/",
			h.GetMe,
			[]core_http_middleware.Middleware{core_http_auth.Middleware(tokenParser)},
		),
	}
}
