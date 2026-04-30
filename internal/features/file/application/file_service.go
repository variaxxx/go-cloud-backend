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

type FileService struct {
	repository file_domain.FileRepository
	storage    file_domain.FileStorage
	folders    folder_domain.FolderRepository
}

func NewFileService(
	repository file_domain.FileRepository,
	storage file_domain.FileStorage,
	folders folder_domain.FolderRepository,
) *FileService {
	return &FileService{
		repository: repository,
		storage:    storage,
		folders:    folders,
	}
}

func (s *FileService) Upload(
	ctx context.Context,
	params UploadFileParams,
) (file_domain.File, error) {
	filename := strings.TrimSpace(params.Filename)
	if filename == "" {
		return file_domain.File{}, fmt.Errorf("upload file: %w", core_errors.ErrInvalidArgument)
	}

	if params.FolderID != nil {
		if _, err := s.folders.FindByIDAndUserID(ctx, *params.FolderID, params.UserID); err != nil {
			return file_domain.File{}, fmt.Errorf("upload file: folder validation: %w", err)
		}
	}

	path, err := s.storage.Save(ctx, filename, params.File)
	if err != nil {
		return file_domain.File{}, fmt.Errorf("file save to storage: %w", err)
	}

	createdFile, err := s.repository.Create(
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

	// TODO: post msg to rmq

	return createdFile, nil
}

func (s *FileService) Delete(
	ctx context.Context,
	id uuid.UUID,
	userID int64,
) error {
	currentFile, err := s.repository.FindByIDAndUserID(ctx, id, userID)
	if err != nil {
		return fmt.Errorf("find current file: %w", err)
	}

	if err := s.repository.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete file: %w", err)
	}

	_ = s.storage.Delete(ctx, currentFile.StoragePath)

	return nil
}
