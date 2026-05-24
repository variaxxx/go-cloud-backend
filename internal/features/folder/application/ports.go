package folder_app

import (
	file_domain "cloud/internal/features/file/domain"
	folder_domain "cloud/internal/features/folder/domain"
	"context"

	"github.com/google/uuid"
)

type FolderRepository interface {
	Create(
		ctx context.Context,
		name string,
		userID int64,
		parentID *uuid.UUID,
	) (folder_domain.Folder, error)

	Update(
		ctx context.Context,
		id uuid.UUID,
		userID int64,
		name string,
		parentID *uuid.UUID,
	) (folder_domain.Folder, error)

	FindByIDAndUserID(
		ctx context.Context,
		id uuid.UUID,
		userID int64,
	) (folder_domain.Folder, error)

	FindByIDAndUserIDWithPath(
		ctx context.Context,
		id uuid.UUID,
		userID int64,
	) (folder_domain.Folder, []folder_domain.Folder, error)

	FindByParentIDAndUserID(
		ctx context.Context,
		parentID *uuid.UUID,
		userID int64,
	) ([]folder_domain.Folder, error)

	WouldCreateCycle(
		ctx context.Context,
		id uuid.UUID,
		parentID uuid.UUID,
		userID int64,
	) (bool, error)

	Delete(
		ctx context.Context,
		id uuid.UUID,
	) error
}

type FileRepository interface {
	FindByFolderIDAndUserID(
		ctx context.Context,
		folderID *uuid.UUID,
		userID int64,
	) ([]file_domain.File, error)

	FindByFolderTreeAndUserID(
		ctx context.Context,
		rootFolderID uuid.UUID,
		userID int64,
	) ([]file_domain.File, error)
}

type FileStorage interface {
	Delete(
		ctx context.Context,
		path string,
	) error
}
