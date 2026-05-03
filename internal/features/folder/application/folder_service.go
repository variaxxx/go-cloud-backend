package folder_app

import (
	core_errors "cloud/internal/core/errors"
	folder_domain "cloud/internal/features/folder/domain"
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type folderService struct {
	folderRepo FolderRepository
	fileRepo   FileRepository
	storage    FileStorage
}

func NewFolderService(
	folderRepo FolderRepository,
	fileRepo FileRepository,
	storage FileStorage,
) *folderService {
	return &folderService{
		folderRepo: folderRepo,
		fileRepo:   fileRepo,
		storage:    storage,
	}
}

func (s *folderService) Create(
	ctx context.Context,
	params CreateFolderParams,
) (folder_domain.Folder, error) {
	name := strings.TrimSpace(params.Name)
	if name == "" {
		return folder_domain.Folder{}, fmt.Errorf("create folder: %w", core_errors.ErrInvalidArgument)
	}

	if params.ParentID != nil {
		if _, err := s.folderRepo.FindByIDAndUserID(ctx, *params.ParentID, params.UserID); err != nil {
			return folder_domain.Folder{}, fmt.Errorf("find parent folder: %w", err)
		}
	}

	folder, err := s.folderRepo.Create(ctx, name, params.UserID, params.ParentID)
	if err != nil {
		return folder_domain.Folder{}, fmt.Errorf("create folder: %w", err)
	}

	return folder, nil
}

func (s *folderService) Edit(
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

func (s *folderService) Delete(
	ctx context.Context,
	id uuid.UUID,
	userID int64,
) error {
	if _, err := s.folderRepo.FindByIDAndUserID(ctx, id, userID); err != nil {
		return fmt.Errorf("find current folder: %w", err)
	}

	files, err := s.fileRepo.FindByFolderTreeAndUserID(ctx, id, userID)
	if err != nil {
		return fmt.Errorf("find files in folder tree: %w", err)
	}

	if err := s.folderRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete folder: %w", err)
	}

	for _, file := range files {
		_ = s.storage.Delete(ctx, file.StoragePath)
	}

	return nil
}

func (s *folderService) GetContent(
	ctx context.Context,
	userID int64,
	folderID *uuid.UUID,
) (FolderContent, error) {
	var folderName *string
	folderPath := make([]string, 0)

	if folderID != nil {
		folder, path, err := s.folderRepo.FindByIDAndUserIDWithPath(ctx, *folderID, userID)
		if err != nil {
			return FolderContent{}, fmt.Errorf("get folder info: %w", err)
		}

		folderName = &folder.Name
		folderPath = path
	}

	nestedFolders, err := s.folderRepo.FindByParentIDAndUserID(ctx, folderID, userID)
	if err != nil {
		return FolderContent{}, fmt.Errorf("get nested folders: %w", err)
	}

	files, err := s.fileRepo.FindByFolderIDAndUserID(ctx, folderID, userID)
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
