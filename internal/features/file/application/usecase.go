package file_app

import (
	file_domain "cloud/internal/features/file/domain"
	"context"
	"io"

	"github.com/google/uuid"
)

type FileUseCase interface {
	Upload(
		ctx context.Context,
		filename string,
		mimetype *string,
		sizeBytes int64,
		userID int64,
		folderID *uuid.UUID,
		file io.Reader,
	) (file_domain.File, error)
}
