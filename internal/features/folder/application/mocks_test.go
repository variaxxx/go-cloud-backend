package folder_app

import (
	file_domain "cloud/internal/features/file/domain"
	folder_domain "cloud/internal/features/folder/domain"
	"context"

	"github.com/google/uuid"
)

type mockFolderRepository struct {
	folder           folder_domain.Folder
	path             []folder_domain.Folder
	nested           []folder_domain.Folder
	err              error
	wouldCreateCycle bool
	findCalls        int
	cycleChecks      int
	updatedParentID  *uuid.UUID
	deletedID        uuid.UUID
}

func (f *mockFolderRepository) Create(ctx context.Context, name string, userID int64, parentID *uuid.UUID) (folder_domain.Folder, error) {
	if f.err != nil {
		return folder_domain.Folder{}, f.err
	}
	return folder_domain.Folder{ID: uuid.New(), Name: name, UserID: userID, ParentID: parentID}, nil
}

func (f *mockFolderRepository) Update(ctx context.Context, id uuid.UUID, userID int64, name string, parentID *uuid.UUID) (folder_domain.Folder, error) {
	if f.err != nil {
		return folder_domain.Folder{}, f.err
	}
	f.updatedParentID = parentID
	f.folder.Name = name
	f.folder.ParentID = parentID
	return f.folder, nil
}

func (f *mockFolderRepository) FindByIDAndUserID(ctx context.Context, id uuid.UUID, userID int64) (folder_domain.Folder, error) {
	f.findCalls++
	if f.err != nil {
		return folder_domain.Folder{}, f.err
	}
	return f.folder, nil
}

func (f *mockFolderRepository) FindByIDAndUserIDWithPath(ctx context.Context, id uuid.UUID, userID int64) (folder_domain.Folder, []folder_domain.Folder, error) {
	if f.err != nil {
		return folder_domain.Folder{}, nil, f.err
	}
	return f.folder, f.path, nil
}

func (f *mockFolderRepository) FindByParentIDAndUserID(ctx context.Context, parentID *uuid.UUID, userID int64) ([]folder_domain.Folder, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.nested, nil
}

func (f *mockFolderRepository) WouldCreateCycle(ctx context.Context, id uuid.UUID, parentID uuid.UUID, userID int64) (bool, error) {
	f.cycleChecks++
	if f.err != nil {
		return false, f.err
	}
	return f.wouldCreateCycle, nil
}

func (f *mockFolderRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if f.err != nil {
		return f.err
	}
	f.deletedID = id
	return nil
}

type mockFolderFileRepository struct {
	files     []file_domain.File
	treeFiles []file_domain.File
	err       error
}

func (f *mockFolderFileRepository) FindByFolderIDAndUserID(ctx context.Context, folderID *uuid.UUID, userID int64) ([]file_domain.File, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.files, nil
}

func (f *mockFolderFileRepository) FindByFolderTreeAndUserID(ctx context.Context, rootFolderID uuid.UUID, userID int64) ([]file_domain.File, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.treeFiles, nil
}

type mockFolderStorage struct {
	deleted []string
	err     error
}

func (f *mockFolderStorage) Delete(ctx context.Context, path string) error {
	if f.err != nil {
		return f.err
	}
	f.deleted = append(f.deleted, path)
	return nil
}
