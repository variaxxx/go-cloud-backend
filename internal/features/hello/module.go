package hello

import (
	core_http_server "cloud/internal/core/transport/http/server"
	hello_app "cloud/internal/features/hello/application"
	hello_http "cloud/internal/features/hello/transport/http"
)

func RegisterHTTP(router *core_http_server.APIRouter) {
	handler := hello_http.NewHelloHTTPHandler(
		hello_app.NewHelloService(),
	)

	router.RegisterRoutes(handler.Routes()...)
}
