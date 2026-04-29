package folder_app

import (
	core_errors "cloud/internal/core/errors"
	folder_domain "cloud/internal/features/folder/domain"
	"context"
	"fmt"
)

type FolderService struct {
	repository folder_domain.FolderRepository
}

func NewFolderService(
	repository folder_domain.FolderRepository,
) *FolderService {
	return &FolderService{
		repository: repository,
	}
}

func (s *FolderService) Create(
	ctx context.Context,
	name string,
	userID int64,
	parentID *int64,
) (folder_domain.Folder, error) {
	if name == "" {
		return folder_domain.Folder{}, fmt.Errorf("create folder: %w", core_errors.ErrInvalidArgument)
	}

	folder, err := s.repository.Create(ctx, name, userID, parentID)
	if err != nil {
		return folder_domain.Folder{}, fmt.Errorf("create folder: %w", err)
	}

	return folder, nil
}
