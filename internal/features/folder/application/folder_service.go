package folder_app

import (
	core_errors "cloud/internal/core/errors"
	file_domain "cloud/internal/features/file/domain"
	folder_domain "cloud/internal/features/folder/domain"
	"context"
	"fmt"
	"strings"
)

type FolderService struct {
	folderRepo folder_domain.FolderRepository
	fileRepo   file_domain.FileRepository
}

func NewFolderService(
	folderRepo folder_domain.FolderRepository,
	fileRepo file_domain.FileRepository,
) *FolderService {
	return &FolderService{
		folderRepo: folderRepo,
		fileRepo:   fileRepo,
	}
}

func (s *FolderService) Create(
	ctx context.Context,
	params CreateFolderParams,
) (folder_domain.Folder, error) {
	if params.Name == "" {
		return folder_domain.Folder{}, fmt.Errorf("create folder: %w", core_errors.ErrInvalidArgument)
	}

	folder, err := s.folderRepo.Create(ctx, params.Name, params.UserID, params.ParentID)
	if err != nil {
		return folder_domain.Folder{}, fmt.Errorf("create folder: %w", err)
	}

	return folder, nil
}

func (s *FolderService) Edit(
	ctx context.Context,
	params EditFolderParams,
) (folder_domain.Folder, error) {
	if params.Name == nil && !params.IsParentIDUpdate {
		return folder_domain.Folder{}, fmt.Errorf("edit folder: %w", core_errors.ErrInvalidArgument)
	}

	currentFolder, err := s.folderRepo.FindByIDAndUserID(ctx, params.ID, params.UserID)
	if err != nil {
		return folder_domain.Folder{}, fmt.Errorf("find current folder: %w", err)
	}

	name := currentFolder.Name
	if params.Name != nil {
		trimmedName := strings.TrimSpace(*params.Name)
		if trimmedName == "" {
			return folder_domain.Folder{}, fmt.Errorf("edit folder: %w", core_errors.ErrInvalidArgument)
		}

		name = trimmedName
	}

	parentID := currentFolder.ParentID
	if params.IsParentIDUpdate {
		parentID = params.ParentID
	}

	if parentID != nil {
		if *parentID == currentFolder.ID {
			return folder_domain.Folder{}, fmt.Errorf("edit folder: %w", core_errors.ErrInvalidArgument)
		}

		if _, err := s.folderRepo.FindByIDAndUserID(ctx, *parentID, params.UserID); err != nil {
			return folder_domain.Folder{}, fmt.Errorf("find parent folder: %w", err)
		}

		wouldCreateCycle, err := s.folderRepo.WouldCreateCycle(ctx, currentFolder.ID, *parentID, params.UserID)
		if err != nil {
			return folder_domain.Folder{}, fmt.Errorf("check folder cycle: %w", err)
		}

		if wouldCreateCycle {
			return folder_domain.Folder{}, fmt.Errorf("edit folder: %w", core_errors.ErrInvalidArgument)
		}
	}

	updatedFolder, err := s.folderRepo.Update(ctx, currentFolder.ID, params.UserID, name, parentID)
	if err != nil {
		return folder_domain.Folder{}, fmt.Errorf("update folder: %w", err)
	}

	return updatedFolder, nil
}

func (s *FolderService) GetContent(
	ctx context.Context,
	params GetFolderContentParams,
) (FolderContent, error) {
	var folderName *string
	folderPath := make([]string, 0)

	if params.FolderID != nil {
		folder, path, err := s.folderRepo.FindByIDAndUserIDWithPath(ctx, *params.FolderID, params.UserID)
		if err != nil {
			return FolderContent{}, fmt.Errorf("get folder info: %w", err)
		}

		folderName = &folder.Name
		folderPath = path
	}

	nestedFolders, err := s.folderRepo.FindByParentIDAndUserID(ctx, params.FolderID, params.UserID)
	if err != nil {
		return FolderContent{}, fmt.Errorf("get nested folders: %w", err)
	}

	files, err := s.fileRepo.FindByFolderIDAndUserID(ctx, params.FolderID, params.UserID)
	if err != nil {
		return FolderContent{}, fmt.Errorf("get files: %w", err)
	}

	return FolderContent{
		FolderName: folderName,
		FolderPath: folderPath,
		Folders:    nestedFolders,
		Files:      files,
	}, nil
}
