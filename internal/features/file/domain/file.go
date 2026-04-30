package file_domain

import (
	"time"

	"github.com/google/uuid"
)

type FileStatus string

const (
	FileStatusUploaded  FileStatus = "uploaded"
	FileStatusProcessed FileStatus = "processed"
	FileStatusFailed    FileStatus = "failed"
)

type File struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
	Filename    string
	Mimetype    *string
	Status      FileStatus
	StoragePath string
	SizeBytes   int64
	UserID      int64
	FolderID    *uuid.UUID
}
