package folder

import (
	core_http_server "cloud/internal/core/transport/http/server"
	auth_jwt "cloud/internal/features/auth/infra/jwt"
	file "cloud/internal/features/file"
	file_postgres "cloud/internal/features/file/infra/postgres"
	file_storage "cloud/internal/features/file/infra/storage"
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
	fileConfig, err := file.NewConfig()
	if err != nil {
		return err
	}

	jwtConfig, err := auth_jwt.NewConfig()
	if err != nil {
		return err
	}

	folderRepo := folder_postgres.NewFolderRepository(deps.DB)
	fileRepo := file_postgres.NewFileRepository(deps.DB)
	fileStorage := file_storage.NewLocalFileStorage(fileConfig.StorageDir)
	folderService := folder_app.NewFolderService(folderRepo, fileRepo, fileStorage)
	tokenManager := auth_jwt.NewManager(jwtConfig)
	handler := folder_http.NewHandler(folderService)

	deps.Router.RegisterRoutes(handler.Routes(tokenManager)...)
	return nil
}
