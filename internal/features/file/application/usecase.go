package file_app

import (
	file_domain "cloud/internal/features/file/domain"
	"context"
	"io"

	"github.com/google/uuid"
)

type UploadFileParams struct {
	Filename  string
	Mimetype  *string
	SizeBytes int64
	UserID    int64
	FolderID  *uuid.UUID
	File      io.Reader
}

type FileUseCase interface {
	Upload(
		ctx context.Context,
		params UploadFileParams,
	) (file_domain.File, error)

	Delete(
		ctx context.Context,
		id uuid.UUID,
		userID int64,
	) error
}
