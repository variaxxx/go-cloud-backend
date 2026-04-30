package folder_http

import (
	core_http_auth "cloud/internal/core/transport/http/auth"
	core_http_middleware "cloud/internal/core/transport/http/middleware"
	core_http_server "cloud/internal/core/transport/http/server"
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
	tokenParser core_http_auth.TokenParser,
) []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/content",
			Handler: h.GetContent,
			Middlewares: []core_http_middleware.Middleware{
				core_http_auth.Middleware(tokenParser),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/{id}/content",
			Handler: h.GetContent,
			Middlewares: []core_http_middleware.Middleware{
				core_http_auth.Middleware(tokenParser),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/",
			Handler: h.Create,
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
	}
}
