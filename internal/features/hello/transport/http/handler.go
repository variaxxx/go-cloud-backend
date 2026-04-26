package hello_http

import (
	core_http_server "cloud/internal/core/transport/http/server"
	hello_app "cloud/internal/features/hello/application"
	"net/http"
)

type HelloHTTPHandler struct {
	service hello_app.HelloUseCase
}

func NewHelloHTTPHandler(
	service hello_app.HelloUseCase,
) *HelloHTTPHandler {
	return &HelloHTTPHandler{
		service: service,
	}
}

func (h *HelloHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/test",
			Handler: h.GetHello,
		},
	}
}
