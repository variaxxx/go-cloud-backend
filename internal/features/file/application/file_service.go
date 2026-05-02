package file_app

import (
	core_errors "cloud/internal/core/errors"
	file_domain "cloud/internal/features/file/domain"
	folder_domain "cloud/internal/features/folder/domain"
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type fileService struct {
	fileRepo       file_domain.FileRepository
	storage        file_domain.FileStorage
	folderRepo     folder_domain.FolderRepository
	eventPublisher FileEventPublisher
}

func NewFileService(
	repository file_domain.FileRepository,
	storage file_domain.FileStorage,
	folders folder_domain.FolderRepository,
	eventPublisher FileEventPublisher,
) *fileService {
	return &fileService{
		fileRepo:       repository,
		storage:        storage,
		folderRepo:     folders,
		eventPublisher: eventPublisher,
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
