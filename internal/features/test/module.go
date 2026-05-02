package test

import (
	core_http_server "cloud/internal/core/transport/http/server"
	test_app "cloud/internal/features/test/application"
	test_postgres "cloud/internal/features/test/infra/postgres"
	test_http "cloud/internal/features/test/transport/http"
	infra_postgres "cloud/internal/infra/postgres"
)

type Deps struct {
	Router *core_http_server.APIRouter
	DB     infra_postgres.Pool
}

func Register(
	deps Deps,
) error {
	repo := test_postgres.NewTestRepository(deps.DB)
	service := test_app.NewTestService(repo)
	handler := test_http.NewHandler(service)

	deps.Router.RegisterRoutes(handler.Routes()...)
	return nil
}
