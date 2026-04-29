package file_http

import (
	core_http_middleware "cloud/internal/core/transport/http/middleware"
	core_http_server "cloud/internal/core/transport/http/server"
	auth_domain "cloud/internal/features/auth/domain"
	auth_http "cloud/internal/features/auth/transport/http"
	file_app "cloud/internal/features/file/application"
	"net/http"
)

type Handler struct {
	fileService file_app.FileUseCase
}

func NewHandler(
	fileService file_app.FileUseCase,
) *Handler {
	return &Handler{
		fileService: fileService,
	}
}

func (h *Handler) Routes(
	tokenManager auth_domain.TokenManager,
) []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/upload",
			Handler: h.Upload,
			Middlewares: []core_http_middleware.Middleware{
				auth_http.Auth(tokenManager),
			},
		},
	}
}
