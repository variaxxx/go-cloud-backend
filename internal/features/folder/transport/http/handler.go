package folder_http

import (
	core_http_auth "cloud/internal/core/transport/http/auth"
	core_http_middleware "cloud/internal/core/transport/http/middleware"
	core_http_server "cloud/internal/core/transport/http/server"
	folder_app "cloud/internal/features/folder/application"
	"net/http"
)

type handler struct {
	folderService folder_app.FolderUseCase
}

func NewHandler(
	folderService folder_app.FolderUseCase,
) *handler {
	return &handler{
		folderService: folderService,
	}
}

func (h *handler) Routes(
	tokenParser core_http_auth.TokenParser,
) []core_http_server.Route {
	return []core_http_server.Route{
		core_http_server.NewRoute(
			http.MethodGet,
			"/content",
			h.GetContent,
			[]core_http_middleware.Middleware{core_http_auth.Middleware(tokenParser)},
		),
		core_http_server.NewRoute(
			http.MethodGet,
			"/{id}/content",
			h.GetContent,
			[]core_http_middleware.Middleware{core_http_auth.Middleware(tokenParser)},
		),
		core_http_server.NewRoute(
			http.MethodPost,
			"/",
			h.Create,
			[]core_http_middleware.Middleware{core_http_auth.Middleware(tokenParser)},
		),
		core_http_server.NewRoute(
			http.MethodPatch,
			"/{id}",
			h.Edit,
			[]core_http_middleware.Middleware{core_http_auth.Middleware(tokenParser)},
		),
		core_http_server.NewRoute(
			http.MethodDelete,
			"/{id}",
			h.Delete,
			[]core_http_middleware.Middleware{core_http_auth.Middleware(tokenParser)},
		),
	}
}
