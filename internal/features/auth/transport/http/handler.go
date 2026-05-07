package auth_http

import (
	core_http_server "cloud/internal/core/transport/http/server"
	auth_app "cloud/internal/features/auth/application"
	"net/http"
	"time"
)

type handler struct {
	authService auth_app.AuthUseCase
	refreshTTL  time.Duration
}

func NewHandler(
	authService auth_app.AuthUseCase,
	refreshTTL time.Duration,
) *handler {
	return &handler{
		authService: authService,
		refreshTTL:  refreshTTL,
	}
}

func (h *handler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		core_http_server.NewRoute(
			http.MethodPost,
			"/register",
			h.Register,
			nil,
		),
		core_http_server.NewRoute(
			http.MethodPost,
			"/login",
			h.Login,
			nil,
		),
		core_http_server.NewRoute(
			http.MethodPost,
			"/refresh",
			h.RefreshTokens,
			nil,
		),
	}
}
