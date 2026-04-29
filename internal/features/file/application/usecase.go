package file_app

import (
	file_domain "cloud/internal/features/file/domain"
	"context"
	"io"
)

type FileUseCase interface {
	Upload(
		ctx context.Context,
		filename string,
		mimetype *string,
		sizeBytes int64,
		userID int64,
		folderID *int64,
		file io.Reader,
	) (file_domain.File, error)
}
