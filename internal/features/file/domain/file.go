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

func (s FileStatus) IsTerminal() bool {
	return s == FileStatusProcessed || s == FileStatusFailed
}

type File struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Filename    string
	Mimetype    *string
	Status      FileStatus
	StoragePath string
	PreviewPath *string
	PreviewMime *string
	SizeBytes   int64
	UserID      int64
	FolderID    *uuid.UUID
}
