package core_http_server

import (
	"fmt"
	"net/http"
)

type APIRouter struct {
	*http.ServeMux
	prefix string
	routes []Route
}

func NewAPIRouter(
	prefix string,
) *APIRouter {
	return &APIRouter{
		ServeMux: http.NewServeMux(),
		prefix:   prefix,
	}
}

func (r *APIRouter) RegisterRoutes(routes ...Route) {
	for _, route := range routes {
		pattern := fmt.Sprintf("%s %s", route.Method, route.Path)
		r.Handle(pattern, route.Handler)
	}

	r.routes = append(r.routes, routes...)
}
