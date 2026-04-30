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
}
