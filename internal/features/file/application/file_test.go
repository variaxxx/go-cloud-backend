package file_app

import (
	core_errors "cloud/internal/core/errors"
	file_domain "cloud/internal/features/file/domain"
	folder_domain "cloud/internal/features/folder/domain"
	"context"
	"errors"
	"image"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestFileServiceUpload(t *testing.T) {
	ctx := context.Background()
	userID := int64(10)
	folderID := uuid.New()
	mime := "text/plain"

	t.Run("trims filename, validates folder, saves file and publishes event", func(t *testing.T) {
		repo := &mockFileRepository{}
		storage := &mockFileStorage{savePath: "objects/readme.txt"}
		folders := &mockFileFolderRepository{folder: folder_domain.Folder{ID: folderID, UserID: userID}}
		publisher := &mockFileEventPublisher{}
		metrics := &mockFileMetrics{}
		service := NewFileService(repo, storage, folders, publisher, metrics)

		got, err := service.Upload(ctx, UploadFileParams{
			Filename:  "  readme.txt  ",
			Mimetype:  &mime,
			SizeBytes: 42,
			UserID:    userID,
			FolderID:  &folderID,
			File:      strings.NewReader("content"),
		})
		if err != nil {
			t.Fatalf("Upload() error = %v", err)
		}

		if got.Filename != "readme.txt" {
			t.Fatalf("created filename = %q, want %q", got.Filename, "readme.txt")
		}
		if folders.findCalls != 1 {
			t.Fatalf("folder validation calls = %d, want 1", folders.findCalls)
		}
		if storage.savedFilename != "readme.txt" {
			t.Fatalf("storage saved filename = %q, want readme.txt", storage.savedFilename)
		}
		if repo.created.StoragePath != "objects/readme.txt" {
			t.Fatalf("repo storage path = %q, want objects/readme.txt", repo.created.StoragePath)
		}
		if len(publisher.uploaded) != 1 {
			t.Fatalf("published events = %d, want 1", len(publisher.uploaded))
		}
		if metrics.uploadedBytes != 42 {
			t.Fatalf("uploaded metric bytes = %d, want 42", metrics.uploadedBytes)
		}
	})

	t.Run("rejects blank filename", func(t *testing.T) {
		service := NewFileService(&mockFileRepository{}, &mockFileStorage{}, &mockFileFolderRepository{}, &mockFileEventPublisher{}, &mockFileMetrics{})

		_, err := service.Upload(ctx, UploadFileParams{Filename: "   ", UserID: userID, File: strings.NewReader("")})
		if !errors.Is(err, core_errors.ErrInvalidArgument) {
			t.Fatalf("Upload() error = %v, want ErrInvalidArgument", err)
		}
	})

	t.Run("deletes stored object when database insert fails", func(t *testing.T) {
		insertErr := errors.New("insert failed")
		repo := &mockFileRepository{createErr: insertErr}
		storage := &mockFileStorage{savePath: "objects/orphan.bin"}
		service := NewFileService(repo, storage, &mockFileFolderRepository{}, &mockFileEventPublisher{}, &mockFileMetrics{})

		_, err := service.Upload(ctx, UploadFileParams{Filename: "orphan.bin", UserID: userID, File: strings.NewReader("data")})
		if !errors.Is(err, insertErr) {
			t.Fatalf("Upload() error = %v, want insert error", err)
		}
		if len(storage.deleted) != 1 || storage.deleted[0] != "objects/orphan.bin" {
			t.Fatalf("deleted storage paths = %v, want [objects/orphan.bin]", storage.deleted)
		}
	})
}

func TestFileServiceEdit(t *testing.T) {
	ctx := context.Background()
	userID := int64(10)
	fileID := uuid.New()
	folderID := uuid.New()
	current := file_domain.File{ID: fileID, Filename: "old.txt", UserID: userID, StoragePath: "old/path"}

	t.Run("updates filename and folder", func(t *testing.T) {
		repo := &mockFileRepository{found: current}
		folders := &mockFileFolderRepository{folder: folder_domain.Folder{ID: folderID, UserID: userID}}
		service := NewFileService(repo, &mockFileStorage{}, folders, &mockFileEventPublisher{}, &mockFileMetrics{})
		name := "  new.txt  "

		got, err := service.Edit(ctx, EditFileParams{
			ID:               fileID,
			UserID:           userID,
			Filename:         &name,
			FolderID:         &folderID,
			IsFolderIDUpdate: true,
		})
		if err != nil {
			t.Fatalf("Edit() error = %v", err)
		}

		if got.Filename != "new.txt" {
			t.Fatalf("updated filename = %q, want new.txt", got.Filename)
		}
		if repo.updatedFolderID == nil || *repo.updatedFolderID != folderID {
			t.Fatalf("updated folderID = %v, want %s", repo.updatedFolderID, folderID)
		}
		if folders.findCalls != 1 {
			t.Fatalf("folder validation calls = %d, want 1", folders.findCalls)
		}
	})

	t.Run("requires at least one update", func(t *testing.T) {
		service := NewFileService(&mockFileRepository{}, &mockFileStorage{}, &mockFileFolderRepository{}, &mockFileEventPublisher{}, &mockFileMetrics{})

		_, err := service.Edit(ctx, EditFileParams{ID: fileID, UserID: userID})
		if !errors.Is(err, core_errors.ErrInvalidArgument) {
			t.Fatalf("Edit() error = %v, want ErrInvalidArgument", err)
		}
	})
}

func TestFileServiceDeleteRemovesDatabaseRecordBeforeStorage(t *testing.T) {
	ctx := context.Background()
	fileID := uuid.New()
	repo := &mockFileRepository{found: file_domain.File{ID: fileID, UserID: 5, StoragePath: "objects/file.bin"}}
	storage := &mockFileStorage{}
	service := NewFileService(repo, storage, &mockFileFolderRepository{}, &mockFileEventPublisher{}, &mockFileMetrics{})

	if err := service.Delete(ctx, fileID, 5); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	if repo.deletedID != fileID {
		t.Fatalf("deleted db id = %s, want %s", repo.deletedID, fileID)
	}
	if len(storage.deleted) != 1 || storage.deleted[0] != "objects/file.bin" {
		t.Fatalf("deleted storage paths = %v, want [objects/file.bin]", storage.deleted)
	}
}

func TestFileServiceDownloadAndPreview(t *testing.T) {
	ctx := context.Background()
	userID := int64(10)
	fileID := uuid.New()
	mime := "image/png"
	previewPath := "previews/file.png"

	t.Run("download opens stored content", func(t *testing.T) {
		repo := &mockFileRepository{found: file_domain.File{ID: fileID, UserID: userID, Filename: "file.png", Mimetype: &mime, StoragePath: "objects/file.png"}}
		storage := &mockFileStorage{openContent: "image-bytes"}
		service := NewFileService(repo, storage, &mockFileFolderRepository{}, &mockFileEventPublisher{}, &mockFileMetrics{})

		got, err := service.Download(ctx, fileID, userID)
		if err != nil {
			t.Fatalf("Download() error = %v", err)
		}
		defer got.Content.Close()

		if got.Filename != "file.png" || got.Mimetype == nil || *got.Mimetype != mime {
			t.Fatalf("Download() metadata = (%q, %v), want file.png/%s", got.Filename, got.Mimetype, mime)
		}
		if storage.openedPath != "objects/file.png" {
			t.Fatalf("opened path = %q, want objects/file.png", storage.openedPath)
		}
	})

	t.Run("preview requires generated preview path", func(t *testing.T) {
		repo := &mockFileRepository{found: file_domain.File{ID: fileID, UserID: userID}}
		service := NewFileService(repo, &mockFileStorage{}, &mockFileFolderRepository{}, &mockFileEventPublisher{}, &mockFileMetrics{})

		_, err := service.GetPreview(ctx, fileID, userID)
		if !errors.Is(err, core_errors.ErrNotFound) {
			t.Fatalf("GetPreview() error = %v, want ErrNotFound", err)
		}
	})

	t.Run("preview opens preview content", func(t *testing.T) {
		repo := &mockFileRepository{found: file_domain.File{ID: fileID, UserID: userID, PreviewPath: &previewPath, PreviewMime: &mime}}
		storage := &mockFileStorage{openContent: "preview"}
		service := NewFileService(repo, storage, &mockFileFolderRepository{}, &mockFileEventPublisher{}, &mockFileMetrics{})

		got, err := service.GetPreview(ctx, fileID, userID)
		if err != nil {
			t.Fatalf("GetPreview() error = %v", err)
		}
		defer got.Content.Close()

		if got.Mimetype == nil || *got.Mimetype != mime {
			t.Fatalf("preview mimetype = %v, want %s", got.Mimetype, mime)
		}
		if storage.openedPath != previewPath {
			t.Fatalf("opened path = %q, want %s", storage.openedPath, previewPath)
		}
	})
}

func TestFileUploadedEventHandler(t *testing.T) {
	ctx := context.Background()
	userID := int64(15)
	fileID := uuid.New()

	t.Run("skips terminal files", func(t *testing.T) {
		repo := &mockFileRepository{found: file_domain.File{
			ID:     fileID,
			UserID: userID,
			Status: file_domain.FileStatusProcessed,
		}}
		metrics := &mockFileMetrics{}
		handler := NewFileUploadedEventHandler(repo, &mockFileStorage{}, metrics)

		if err := handler.HandleUploaded(ctx, FileUploadedEvent{FileID: fileID, UserID: userID}); err != nil {
			t.Fatalf("HandleUploaded() error = %v", err)
		}
		if metrics.started != 0 {
			t.Fatalf("processing started = %d, want 0", metrics.started)
		}
	})

	t.Run("marks unsupported files as processed", func(t *testing.T) {
		repo := &mockFileRepository{found: file_domain.File{
			ID:       fileID,
			UserID:   userID,
			Filename: "notes.txt",
			Status:   file_domain.FileStatusUploaded,
		}}
		metrics := &mockFileMetrics{}
		handler := NewFileUploadedEventHandler(repo, &mockFileStorage{}, metrics)

		if err := handler.HandleUploaded(ctx, FileUploadedEvent{FileID: fileID, UserID: userID}); err != nil {
			t.Fatalf("HandleUploaded() error = %v", err)
		}
		if repo.found.Status != file_domain.FileStatusProcessed {
			t.Fatalf("status = %s, want processed", repo.found.Status)
		}
		if len(metrics.finished) != 1 || metrics.finished[0] != FileProcessingResultSuccess {
			t.Fatalf("finished metrics = %v, want success", metrics.finished)
		}
	})

	t.Run("creates image preview and marks file processed", func(t *testing.T) {
		mime := "image/png"
		storagePath := "objects/picture.png"
		repo := &mockFileRepository{found: file_domain.File{
			ID:          fileID,
			UserID:      userID,
			Filename:    "picture.png",
			Mimetype:    &mime,
			Status:      file_domain.FileStatusUploaded,
			StoragePath: storagePath,
		}}
		storage := &mockFileStorage{openBytes: testJPEG(t, 640, 320)}
		metrics := &mockFileMetrics{}
		handler := NewFileUploadedEventHandler(repo, storage, metrics)

		if err := handler.HandleUploaded(ctx, FileUploadedEvent{FileID: fileID, UserID: userID}); err != nil {
			t.Fatalf("HandleUploaded() error = %v", err)
		}

		if repo.found.Status != file_domain.FileStatusProcessed {
			t.Fatalf("status = %s, want processed", repo.found.Status)
		}
		if repo.found.PreviewPath == nil || *repo.found.PreviewPath != "previews/objects/picture.png.jpg" {
			t.Fatalf("preview path = %v, want previews/objects/picture.png.jpg", repo.found.PreviewPath)
		}
		if storage.savedAtPath != "previews/objects/picture.png.jpg" {
			t.Fatalf("saved preview path = %q, want previews/objects/picture.png.jpg", storage.savedAtPath)
		}
		if storage.savedAtBytes == 0 {
			t.Fatal("saved preview is empty")
		}
		if len(metrics.finished) != 1 || metrics.finished[0] != FileProcessingResultSuccess {
			t.Fatalf("finished metrics = %v, want success", metrics.finished)
		}
	})

	t.Run("marks file failed when preview generation fails", func(t *testing.T) {
		mime := "image/png"
		repo := &mockFileRepository{found: file_domain.File{
			ID:          fileID,
			UserID:      userID,
			Filename:    "broken.png",
			Mimetype:    &mime,
			Status:      file_domain.FileStatusUploaded,
			StoragePath: "objects/broken.png",
		}}
		storage := &mockFileStorage{openErr: errors.New("open failed")}
		metrics := &mockFileMetrics{}
		handler := NewFileUploadedEventHandler(repo, storage, metrics)

		err := handler.HandleUploaded(ctx, FileUploadedEvent{FileID: fileID, UserID: userID})
		if err == nil {
			t.Fatal("HandleUploaded() error = nil, want preview error")
		}
		if repo.found.Status != file_domain.FileStatusFailed {
			t.Fatalf("status = %s, want failed", repo.found.Status)
		}
		if len(metrics.finished) != 1 || metrics.finished[0] != FileProcessingResultFailed {
			t.Fatalf("finished metrics = %v, want failed", metrics.finished)
		}
	})
}

func TestPreviewHelpers(t *testing.T) {
	pdfMime := "application/pdf"
	imageMime := "IMAGE/PNG"

	tests := []struct {
		name string
		file file_domain.File
		want bool
	}{
		{name: "image mimetype", file: file_domain.File{Filename: "x.bin", Mimetype: &imageMime}, want: true},
		{name: "pdf mimetype", file: file_domain.File{Filename: "x.bin", Mimetype: &pdfMime}, want: true},
		{name: "image extension", file: file_domain.File{Filename: "photo.jpeg"}, want: true},
		{name: "unsupported", file: file_domain.File{Filename: "notes.txt"}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isPreviewSupported(tt.file); got != tt.want {
				t.Fatalf("isPreviewSupported() = %v, want %v", got, tt.want)
			}
		})
	}

	resized := resizeImage(image.NewRGBA(image.Rect(0, 0, 640, 320)), 320)
	if resized.Bounds().Dx() != 320 || resized.Bounds().Dy() != 160 {
		t.Fatalf("resized bounds = %v, want 320x160", resized.Bounds())
	}
}
