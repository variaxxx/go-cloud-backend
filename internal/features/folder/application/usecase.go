package folder_app

import (
	file_domain "cloud/internal/features/file/domain"
	folder_domain "cloud/internal/features/folder/domain"
	"context"

	"github.com/google/uuid"
)

type CreateFolderParams struct {
	Name     string
	UserID   int64
	ParentID *uuid.UUID
}

type EditFolderParams struct {
	ID               uuid.UUID
	UserID           int64
	Name             *string
	ParentID         *uuid.UUID
	IsParentIDUpdate bool
}

type GetFolderContentParams struct {
	UserID   int64
	FolderID *uuid.UUID
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

	Edit(
		ctx context.Context,
		params EditFolderParams,
	) (folder_domain.Folder, error)

	Delete(
		ctx context.Context,
		id uuid.UUID,
		userID int64,
	) error

	GetContent(
		ctx context.Context,
		params GetFolderContentParams,
	) (FolderContent, error)
}
