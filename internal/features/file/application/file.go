package file_app

import (
	core_errors "cloud/internal/core/errors"
	file_domain "cloud/internal/features/file/domain"
	"context"
	"fmt"
	"io"
	"strings"

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

type EditFileParams struct {
	ID               uuid.UUID
	UserID           int64
	Filename         *string
	FolderID         *uuid.UUID
	IsFolderIDUpdate bool
}

type DownloadFileResult struct {
	Filename string
	Mimetype *string
	Content  io.ReadCloser
}

type FilePreview struct {
	Mimetype *string
	Content  io.ReadCloser
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

	Edit(
		ctx context.Context,
		params EditFileParams,
	) (file_domain.File, error)

	Download(
		ctx context.Context,
		id uuid.UUID,
		userID int64,
	) (DownloadFileResult, error)

	GetPreview(
		ctx context.Context,
		id uuid.UUID,
		userID int64,
	) (FilePreview, error)
}

type fileService struct {
	fileRepo       FileRepository
	storage        FileStorage
	folderRepo     FolderRepository
	eventPublisher FileEventPublisher
	metrics        FileMetrics
}

func NewFileService(
	repository FileRepository,
	storage FileStorage,
	folders FolderRepository,
	eventPublisher FileEventPublisher,
	metrics FileMetrics,
) *fileService {
	return &fileService{
		fileRepo:       repository,
		storage:        storage,
		folderRepo:     folders,
		eventPublisher: eventPublisher,
		metrics:        metrics,
	}
}

func (s *fileService) Upload(
	ctx context.Context,
	params UploadFileParams,
) (file_domain.File, error) {
	filename := strings.TrimSpace(params.Filename)
	if filename == "" {
		return file_domain.File{}, fmt.Errorf("upload file: %w", core_errors.ErrInvalidArgument)
	}

	if params.FolderID != nil {
		if _, err := s.folderRepo.FindByIDAndUserID(ctx, *params.FolderID, params.UserID); err != nil {
			return file_domain.File{}, fmt.Errorf("upload file: folder validation: %w", err)
		}
	}

	path, err := s.storage.Save(ctx, filename, params.File)
	if err != nil {
		return file_domain.File{}, fmt.Errorf("file save to storage: %w", err)
	}

	createdFile, err := s.fileRepo.Create(
		ctx,
		filename,
		params.Mimetype,
		file_domain.FileStatusUploaded,
		path,
		params.SizeBytes,
		params.UserID,
		params.FolderID,
	)
	if err != nil {
		if deleteErr := s.storage.Delete(ctx, path); deleteErr != nil {
			return file_domain.File{}, fmt.Errorf("create file record: %w; cleanup stored file: %v", err, deleteErr)
		}
		return file_domain.File{}, fmt.Errorf("create file record: %w", err)
	}

	_ = s.eventPublisher.PublishUploaded(ctx, createdFile)
	s.metrics.UploadSucceeded(params.SizeBytes)

	return createdFile, nil
}

func (s *fileService) Delete(
	ctx context.Context,
	id uuid.UUID,
	userID int64,
) error {
	currentFile, err := s.fileRepo.FindByIDAndUserID(ctx, id, userID)
	if err != nil {
		return fmt.Errorf("find current file: %w", err)
	}

	if err := s.fileRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete file: %w", err)
	}

	_ = s.storage.Delete(ctx, currentFile.StoragePath)

	return nil
}

func (s *fileService) Edit(
	ctx context.Context,
	params EditFileParams,
) (file_domain.File, error) {
	if params.Filename == nil && !params.IsFolderIDUpdate {
		return file_domain.File{}, fmt.Errorf("edit file: %w", core_errors.ErrInvalidArgument)
	}

	currentFile, err := s.fileRepo.FindByIDAndUserID(ctx, params.ID, params.UserID)
	if err != nil {
		return file_domain.File{}, fmt.Errorf("find current file: %w", err)
	}

	filename := currentFile.Filename
	if params.Filename != nil {
		trimmedName := strings.TrimSpace(*params.Filename)
		if trimmedName == "" {
			return file_domain.File{}, fmt.Errorf("edit file: %w", core_errors.ErrInvalidArgument)
		}

		filename = trimmedName
	}

	folderID := currentFile.FolderID
	if params.IsFolderIDUpdate {
		folderID = params.FolderID
	}

	if folderID != nil {
		if _, err := s.folderRepo.FindByIDAndUserID(ctx, *folderID, params.UserID); err != nil {
			return file_domain.File{}, fmt.Errorf("find folder: %w", err)
		}
	}

	updatedFile, err := s.fileRepo.Update(ctx, currentFile.ID, params.UserID, filename, folderID)
	if err != nil {
		return file_domain.File{}, fmt.Errorf("update file: %w", err)
	}

	return updatedFile, nil
}

func (s *fileService) Download(
	ctx context.Context,
	id uuid.UUID,
	userID int64,
) (DownloadFileResult, error) {
	file, err := s.fileRepo.FindByIDAndUserID(ctx, id, userID)
	if err != nil {
		return DownloadFileResult{}, fmt.Errorf("find current file: %w", err)
	}

	content, err := s.storage.Open(ctx, file.StoragePath)
	if err != nil {
		return DownloadFileResult{}, fmt.Errorf("open file content: %w", err)
	}

	return DownloadFileResult{
		Filename: file.Filename,
		Mimetype: file.Mimetype,
		Content:  content,
	}, nil
}

func (s *fileService) GetPreview(
	ctx context.Context,
	id uuid.UUID,
	userID int64,
) (FilePreview, error) {
	file, err := s.fileRepo.FindByIDAndUserID(ctx, id, userID)
	if err != nil {
		return FilePreview{}, fmt.Errorf("find current file: %w", err)
	}

	if file.PreviewPath == nil {
		return FilePreview{}, fmt.Errorf("find file preview: %w", core_errors.ErrNotFound)
	}

	content, err := s.storage.Open(ctx, *file.PreviewPath)
	if err != nil {
		return FilePreview{}, fmt.Errorf("open file preview: %w", err)
	}

	return FilePreview{
		Mimetype: file.PreviewMime,
		Content:  content,
	}, nil
}
