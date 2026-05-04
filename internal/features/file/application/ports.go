package file_app

import (
	file_domain "cloud/internal/features/file/domain"
	folder_domain "cloud/internal/features/folder/domain"
	"context"
	"io"
	"time"

	"github.com/google/uuid"
)

type FileRepository interface {
	Create(
		ctx context.Context,
		filename string,
		mimetype *string,
		status file_domain.FileStatus,
		storagePath string,
		sizeBytes int64,
		userID int64,
		folderID *uuid.UUID,
	) (file_domain.File, error)

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

	FindByIDAndUserID(
		ctx context.Context,
		id uuid.UUID,
		userID int64,
	) (file_domain.File, error)

	Update(
		ctx context.Context,
		id uuid.UUID,
		userID int64,
		filename string,
		folderID *uuid.UUID,
	) (file_domain.File, error)

	UpdateStatus(
		ctx context.Context,
		id uuid.UUID,
		userID int64,
		status file_domain.FileStatus,
	) (file_domain.File, error)

	UpdatePreview(
		ctx context.Context,
		id uuid.UUID,
		userID int64,
		status file_domain.FileStatus,
		previewPath *string,
		previewMimetype *string,
	) (file_domain.File, error)

	Delete(
		ctx context.Context,
		id uuid.UUID,
	) error
}

type FileStorage interface {
	Save(
		ctx context.Context,
		filename string,
		src io.Reader,
	) (path string, err error)

	SaveAtPath(
		ctx context.Context,
		path string,
		src io.Reader,
	) error

	Delete(
		ctx context.Context,
		path string,
	) error

	Open(
		ctx context.Context,
		path string,
	) (io.ReadCloser, error)
}

type FolderRepository interface {
	FindByIDAndUserID(
		ctx context.Context,
		id uuid.UUID,
		userID int64,
	) (folder_domain.Folder, error)
}

type FileMetrics interface {
	UploadSucceeded(sizeBytes int64)
	ProcessingStarted()
	ObserveProcessingFinished(result FileProcessingResult, duration time.Duration)
}

type FileProcessingResult string

const (
	FileProcessingResultSuccess FileProcessingResult = "success"
	FileProcessingResultFailed  FileProcessingResult = "failed"
)

type FilePreviewResult string

const (
	FilePreviewResultGenerated FilePreviewResult = "generated"
	FilePreviewResultFailed    FilePreviewResult = "failed"
)
