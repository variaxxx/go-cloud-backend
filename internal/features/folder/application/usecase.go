package folder_app

import (
	file_domain "cloud/internal/features/file/domain"
	folder_domain "cloud/internal/features/folder/domain"
	"context"
)

type CreateFolderParams struct {
	Name     string
	UserID   int64
	ParentID *int64
}

type GetFolderContentParams struct {
	UserID   int64
	FolderID *int64
}

type FolderContent struct {
	FolderName *string
	FolderPath []string
	Folders    []folder_domain.Folder
	Files      []file_domain.File
}

type FolderUseCase interface {
	Create(
		ctx context.Context,
		params CreateFolderParams,
	) (folder_domain.Folder, error)

	GetContent(
		ctx context.Context,
		params GetFolderContentParams,
	) (FolderContent, error)
}
