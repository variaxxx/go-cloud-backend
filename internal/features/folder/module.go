package folder

import (
	core_http_server "cloud/internal/core/transport/http/server"
	auth_jwt "cloud/internal/features/auth/infra/jwt"
	file_postgres "cloud/internal/features/file/infra/postgres"
	folder_app "cloud/internal/features/folder/application"
	folder_postgres "cloud/internal/features/folder/infra/postgres"
	folder_http "cloud/internal/features/folder/transport/http"
	infra_postgres "cloud/internal/infra/postgres"
)

type Deps struct {
	Router *core_http_server.APIRouter
	DB     infra_postgres.Pool
}

func Register(
	deps Deps,
) error {
	jwtConfig, err := auth_jwt.NewConfig()
	if err != nil {
		return err
	}

	folderRepo := folder_postgres.NewFolderRepository(deps.DB)
	fileRepo := file_postgres.NewFileRepository(deps.DB)
	service := folder_app.NewFolderService(folderRepo, fileRepo)
	tokenManager := auth_jwt.NewManager(jwtConfig)
	handler := folder_http.NewHandler(service)

	deps.Router.RegisterRoutes(handler.Routes(tokenManager)...)
	return nil
}
