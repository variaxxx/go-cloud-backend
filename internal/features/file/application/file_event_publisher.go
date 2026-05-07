package file_app

import (
	file_domain "cloud/internal/features/file/domain"
	"context"

	"github.com/google/uuid"
)

type FileUploadedEvent struct {
	FileID      uuid.UUID `json:"file_id"`
	Filename    string    `json:"filename"`
	Mimetype    *string
	StoragePath string
	UserID      int64 `json:"user_id"`
}

type FileEventPublisher interface {
	PublishUploaded(
		ctx context.Context,
		file file_domain.File,
	) error
}
