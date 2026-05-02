package file_domain

import (
	"context"

	"github.com/google/uuid"
)

type FileRepository interface {
	Create(
		ctx context.Context,
		filename string,
		mimetype *string,
		status FileStatus,
		storagePath string,
		sizeBytes int64,
		userID int64,
		folderID *uuid.UUID,
	) (File, error)

	FindByFolderIDAndUserID(
		ctx context.Context,
		folderID *uuid.UUID,
		userID int64,
	) ([]File, error)

	FindByFolderTreeAndUserID(
		ctx context.Context,
		rootFolderID uuid.UUID,
		userID int64,
	) ([]File, error)

	FindByIDAndUserID(
		ctx context.Context,
		id uuid.UUID,
		userID int64,
	) (File, error)

	Update(
		ctx context.Context,
		id uuid.UUID,
		userID int64,
		filename string,
		folderID *uuid.UUID,
	) (File, error)

	UpdateStatus(
		ctx context.Context,
		id uuid.UUID,
		userID int64,
		status FileStatus,
	) (File, error)

	UpdatePreview(
		ctx context.Context,
		id uuid.UUID,
		userID int64,
		status FileStatus,
		previewPath *string,
		previewMimetype *string,
	) (File, error)

	Delete(
		ctx context.Context,
		id uuid.UUID,
	) error
}
