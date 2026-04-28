package auth_http

import (
	core_http_server "cloud/internal/core/transport/http/server"
	auth_app "cloud/internal/features/auth/application"
	"net/http"
	"time"
)

type Handler struct {
	authService auth_app.AuthUseCase
	refreshTTL  time.Duration
}

func NewHandler(
	authService auth_app.AuthUseCase,
	refreshTTL time.Duration,
) *Handler {
	return &Handler{
		authService: authService,
		refreshTTL:  refreshTTL,
	}
}

func (h *Handler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/register",
			Handler: h.Register,
		},
		{
			Method:  http.MethodPost,
			Path:    "/login",
			Handler: h.Login,
		},
	}
}
