package file_app

import (
	core_errors "cloud/internal/core/errors"
	file_domain "cloud/internal/features/file/domain"
	folder_domain "cloud/internal/features/folder/domain"
	"context"
	"fmt"
	"io"
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
	filename string,
	mimetype *string,
	sizeBytes int64,
	userID int64,
	folderID *uuid.UUID,
	file io.Reader,
) (file_domain.File, error) {
	filename = strings.TrimSpace(filename)
	if filename == "" {
		return file_domain.File{}, fmt.Errorf("upload file: %w", core_errors.ErrInvalidArgument)
	}

	if folderID != nil {
		if _, err := s.folders.FindByIDAndUserID(ctx, *folderID, userID); err != nil {
			return file_domain.File{}, fmt.Errorf("upload file: folder validation: %w", err)
		}
	}

	path, err := s.storage.Save(ctx, filename, file)
	if err != nil {
		return file_domain.File{}, fmt.Errorf("file save to storage: %w", err)
	}

	createdFile, err := s.repository.Create(
		ctx,
		filename,
		mimetype,
		file_domain.FileStatusUploaded,
		path,
		sizeBytes,
		userID,
		folderID,
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
