package hello_http

import (
	core_http_server "cloud/internal/core/transport/http/server"
	hello_app "cloud/internal/features/hello/application"
	"net/http"
)

type helloHTTPHandler struct {
	service hello_app.HelloUseCase
}

func NewHelloHTTPHandler(
	service hello_app.HelloUseCase,
) *helloHTTPHandler {
	return &helloHTTPHandler{
		service: service,
	}
}

func (h *helloHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		core_http_server.NewRoute(
			http.MethodGet,
			"/test",
			h.GetHello,
			nil,
		),
	}
}
