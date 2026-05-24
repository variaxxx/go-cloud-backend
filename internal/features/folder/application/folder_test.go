package folder_app

import (
	core_errors "cloud/internal/core/errors"
	file_domain "cloud/internal/features/file/domain"
	folder_domain "cloud/internal/features/folder/domain"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
)

func TestFolderServiceCreate(t *testing.T) {
	ctx := context.Background()
	userID := int64(7)
	parentID := uuid.New()

	t.Run("trims name and validates parent", func(t *testing.T) {
		repo := &mockFolderRepository{folder: folder_domain.Folder{ID: parentID, UserID: userID}}
		service := NewFolderService(repo, &mockFolderFileRepository{}, &mockFolderStorage{})

		got, err := service.Create(ctx, CreateFolderParams{Name: "  docs  ", UserID: userID, ParentID: &parentID})
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		if got.Name != "docs" {
			t.Fatalf("created name = %q, want docs", got.Name)
		}
		if repo.findCalls != 1 {
			t.Fatalf("parent validation calls = %d, want 1", repo.findCalls)
		}
	})

	t.Run("rejects blank name", func(t *testing.T) {
		service := NewFolderService(&mockFolderRepository{}, &mockFolderFileRepository{}, &mockFolderStorage{})

		_, err := service.Create(ctx, CreateFolderParams{Name: "   ", UserID: userID})
		if !errors.Is(err, core_errors.ErrInvalidArgument) {
			t.Fatalf("Create() error = %v, want ErrInvalidArgument", err)
		}
	})
}

func TestFolderServiceEdit(t *testing.T) {
	ctx := context.Background()
	userID := int64(7)
	folderID := uuid.New()
	parentID := uuid.New()
	current := folder_domain.Folder{ID: folderID, Name: "old", UserID: userID}

	t.Run("updates name and parent", func(t *testing.T) {
		repo := &mockFolderRepository{folder: current}
		service := NewFolderService(repo, &mockFolderFileRepository{}, &mockFolderStorage{})
		name := "  new  "

		got, err := service.Edit(ctx, EditFolderParams{
			ID:               folderID,
			UserID:           userID,
			Name:             &name,
			ParentID:         &parentID,
			IsParentIDUpdate: true,
		})
		if err != nil {
			t.Fatalf("Edit() error = %v", err)
		}

		if got.Name != "new" {
			t.Fatalf("updated name = %q, want new", got.Name)
		}
		if repo.updatedParentID == nil || *repo.updatedParentID != parentID {
			t.Fatalf("updated parentID = %v, want %s", repo.updatedParentID, parentID)
		}
		if repo.cycleChecks != 1 {
			t.Fatalf("cycle checks = %d, want 1", repo.cycleChecks)
		}
	})

	t.Run("rejects no-op edit", func(t *testing.T) {
		service := NewFolderService(&mockFolderRepository{}, &mockFolderFileRepository{}, &mockFolderStorage{})

		_, err := service.Edit(ctx, EditFolderParams{ID: folderID, UserID: userID})
		if !errors.Is(err, core_errors.ErrInvalidArgument) {
			t.Fatalf("Edit() error = %v, want ErrInvalidArgument", err)
		}
	})

	t.Run("rejects self parent", func(t *testing.T) {
		repo := &mockFolderRepository{folder: current}
		service := NewFolderService(repo, &mockFolderFileRepository{}, &mockFolderStorage{})

		_, err := service.Edit(ctx, EditFolderParams{ID: folderID, UserID: userID, ParentID: &folderID, IsParentIDUpdate: true})
		if !errors.Is(err, core_errors.ErrInvalidArgument) {
			t.Fatalf("Edit() error = %v, want ErrInvalidArgument", err)
		}
	})

	t.Run("rejects cycle", func(t *testing.T) {
		repo := &mockFolderRepository{folder: current, wouldCreateCycle: true}
		service := NewFolderService(repo, &mockFolderFileRepository{}, &mockFolderStorage{})

		_, err := service.Edit(ctx, EditFolderParams{ID: folderID, UserID: userID, ParentID: &parentID, IsParentIDUpdate: true})
		if !errors.Is(err, core_errors.ErrInvalidArgument) {
			t.Fatalf("Edit() error = %v, want ErrInvalidArgument", err)
		}
	})
}

func TestFolderServiceDeleteCleansStoredFilesAfterDatabaseDelete(t *testing.T) {
	ctx := context.Background()
	userID := int64(7)
	folderID := uuid.New()
	repo := &mockFolderRepository{folder: folder_domain.Folder{ID: folderID, UserID: userID}}
	files := &mockFolderFileRepository{treeFiles: []file_domain.File{
		{ID: uuid.New(), StoragePath: "objects/a.txt"},
		{ID: uuid.New(), StoragePath: "objects/b.txt"},
	}}
	storage := &mockFolderStorage{}
	service := NewFolderService(repo, files, storage)

	if err := service.Delete(ctx, folderID, userID); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	if repo.deletedID != folderID {
		t.Fatalf("deleted folder id = %s, want %s", repo.deletedID, folderID)
	}
	if len(storage.deleted) != 2 {
		t.Fatalf("deleted storage paths = %v, want 2 paths", storage.deleted)
	}
}

func TestFolderServiceGetContent(t *testing.T) {
	ctx := context.Background()
	userID := int64(7)
	folderID := uuid.New()
	folderName := "docs"
	nested := []folder_domain.Folder{{ID: uuid.New(), Name: "nested", UserID: userID}}
	files := []file_domain.File{{ID: uuid.New(), Filename: "a.txt", UserID: userID}}
	rootID := uuid.New()
	repo := &mockFolderRepository{
		folder: folder_domain.Folder{ID: folderID, Name: folderName, UserID: userID},
		path: []folder_domain.Folder{
			{ID: rootID, Name: "root", UserID: userID},
			{ID: folderID, Name: "docs", UserID: userID},
		},
		nested: nested,
	}
	fileRepo := &mockFolderFileRepository{files: files}
	service := NewFolderService(repo, fileRepo, &mockFolderStorage{})

	got, err := service.GetContent(ctx, userID, &folderID)
	if err != nil {
		t.Fatalf("GetContent() error = %v", err)
	}

	if got.FolderName == nil || *got.FolderName != folderName {
		t.Fatalf("folder name = %v, want %s", got.FolderName, folderName)
	}
	if len(got.FolderPath) != 2 || got.FolderPath[0].ID != rootID || got.FolderPath[1].ID != folderID || got.FolderPath[1].Name != "docs" {
		t.Fatalf("folder path = %v, want root and docs folders", got.FolderPath)
	}
	if len(got.Folders) != 1 || got.Folders[0].ID != nested[0].ID {
		t.Fatalf("folders = %v, want nested folder", got.Folders)
	}
	if len(got.Files) != 1 || got.Files[0].ID != files[0].ID {
		t.Fatalf("files = %v, want file", got.Files)
	}
}
