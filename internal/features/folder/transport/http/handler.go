package folder_http

import (
	core_http_middleware "cloud/internal/core/transport/http/middleware"
	core_http_server "cloud/internal/core/transport/http/server"
	auth_domain "cloud/internal/features/auth/domain"
	auth_http "cloud/internal/features/auth/transport/http"
	folder_app "cloud/internal/features/folder/application"
	"net/http"
)

type Handler struct {
	folderService folder_app.FolderUseCase
}

func NewHandler(
	folderService folder_app.FolderUseCase,
) *Handler {
	return &Handler{
		folderService: folderService,
	}
}

func (h *Handler) Routes(
	tokenManager auth_domain.TokenManager,
) []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/",
			Handler: h.Create,
			Middlewares: []core_http_middleware.Middleware{
				auth_http.Auth(tokenManager),
			},
		},
	}
}
