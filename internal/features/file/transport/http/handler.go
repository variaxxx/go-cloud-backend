package file_http

import (
	core_http_auth "cloud/internal/core/transport/http/auth"
	core_http_middleware "cloud/internal/core/transport/http/middleware"
	core_http_server "cloud/internal/core/transport/http/server"
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
	tokenParser core_http_auth.TokenParser,
) []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/upload",
			Handler: h.Upload,
			Middlewares: []core_http_middleware.Middleware{
				core_http_auth.Middleware(tokenParser),
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/{id}",
			Handler: h.Edit,
			Middlewares: []core_http_middleware.Middleware{
				core_http_auth.Middleware(tokenParser),
			},
		},
		{
			Method:  http.MethodDelete,
			Path:    "/{id}",
			Handler: h.Delete,
			Middlewares: []core_http_middleware.Middleware{
				core_http_auth.Middleware(tokenParser),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/{id}/download",
			Handler: h.Download,
			Middlewares: []core_http_middleware.Middleware{
				core_http_auth.Middleware(tokenParser),
			},
		},
	}
}
