package file

import (
	core_http_server "cloud/internal/core/transport/http/server"
	auth_jwt "cloud/internal/features/auth/infra/jwt"
	file_app "cloud/internal/features/file/application"
	file_postgres "cloud/internal/features/file/infra/postgres"
	file_storage "cloud/internal/features/file/infra/storage"
	file_http "cloud/internal/features/file/transport/http"
	folder_postgres "cloud/internal/features/folder/infra/postgres"
	infra_postgres "cloud/internal/infra/postgres"
)

type Deps struct {
	Router *core_http_server.APIRouter
	DB     infra_postgres.Pool
}

func Register(
	deps Deps,
) error {
	config, err := NewConfig()
	if err != nil {
		return err
	}

	jwtConfig, err := auth_jwt.NewConfig()
	if err != nil {
		return err
	}

	fileRepo := file_postgres.NewFileRepository(deps.DB)
	folderRepo := folder_postgres.NewFolderRepository(deps.DB)
	storage := file_storage.NewLocalFileStorage(config.StorageDir)
	tokenManager := auth_jwt.NewManager(jwtConfig)

	service := file_app.NewFileService(fileRepo, storage, folderRepo)

	handler := file_http.NewHandler(service)

	deps.Router.RegisterRoutes(handler.Routes(tokenManager)...)
	return nil
}
