package file_http

import (
	core_http_auth "cloud/internal/core/transport/http/auth"
	core_http_middleware "cloud/internal/core/transport/http/middleware"
	core_http_server "cloud/internal/core/transport/http/server"
	file_app "cloud/internal/features/file/application"
	"net/http"
)

type handler struct {
	fileService file_app.FileUseCase
}

func NewHandler(
	fileService file_app.FileUseCase,
) *handler {
	return &handler{
		fileService: fileService,
	}
}

func (h *handler) Routes(
	tokenParser core_http_auth.TokenParser,
) []core_http_server.Route {
	return []core_http_server.Route{
		core_http_server.NewRoute(
			http.MethodPost,
			"/upload",
			h.Upload,
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
		core_http_server.NewRoute(
			http.MethodGet,
			"/{id}/download",
			h.Download,
			[]core_http_middleware.Middleware{core_http_auth.Middleware(tokenParser)},
		),
	}
}
