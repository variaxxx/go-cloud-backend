package core_http_server

import (
	core_http_middleware "cloud/internal/core/transport/http/middleware"
	"net/http"
)

type Route struct {
	Method      string
	Path        string
	Handler     http.HandlerFunc
	Middlewares []core_http_middleware.Middleware
}

func NewRoute(
	method string,
	path string,
	handler http.HandlerFunc,
	middlewares []core_http_middleware.Middleware,
) Route {
	return Route{
		Method:      method,
		Path:        path,
		Handler:     handler,
		Middlewares: middlewares,
	}
}
