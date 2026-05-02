package test_http

import (
	core_http_server "cloud/internal/core/transport/http/server"
	test_app "cloud/internal/features/test/application"
	"net/http"
)

type Handler struct {
	service test_app.TestUseCase
}

func NewHandler(
	service test_app.TestUseCase,
) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		core_http_server.NewRoute(
			http.MethodGet,
			"/test",
			h.GetHello,
			nil,
		),
		core_http_server.NewRoute(
			http.MethodPost,
			"/dbtest",
			h.DbTest,
			nil,
		),
	}
}
