package file_app

import (
	"bytes"
	file_domain "cloud/internal/features/file/domain"
	folder_domain "cloud/internal/features/folder/domain"
	"context"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

type mockFileRepository struct {
	createErr       error
	findErr         error
	updateErr       error
	deleteErr       error
	found           file_domain.File
	created         file_domain.File
	updatedFolderID *uuid.UUID
	deletedID       uuid.UUID
}

func (f *mockFileRepository) Create(ctx context.Context, filename string, mimetype *string, status file_domain.FileStatus, storagePath string, sizeBytes int64, userID int64, folderID *uuid.UUID) (file_domain.File, error) {
	if f.createErr != nil {
		return file_domain.File{}, f.createErr
	}
	f.created = file_domain.File{
		ID:          uuid.New(),
		Filename:    filename,
		Mimetype:    mimetype,
		Status:      status,
		StoragePath: storagePath,
		SizeBytes:   sizeBytes,
		UserID:      userID,
		FolderID:    folderID,
	}
	return f.created, nil
}

func (f *mockFileRepository) FindByFolderIDAndUserID(ctx context.Context, folderID *uuid.UUID, userID int64) ([]file_domain.File, error) {
	return nil, nil
}

func (f *mockFileRepository) FindByFolderTreeAndUserID(ctx context.Context, rootFolderID uuid.UUID, userID int64) ([]file_domain.File, error) {
	return nil, nil
}

func (f *mockFileRepository) FindByIDAndUserID(ctx context.Context, id uuid.UUID, userID int64) (file_domain.File, error) {
	if f.findErr != nil {
		return file_domain.File{}, f.findErr
	}
	return f.found, nil
}

func (f *mockFileRepository) Update(ctx context.Context, id uuid.UUID, userID int64, filename string, folderID *uuid.UUID) (file_domain.File, error) {
	if f.updateErr != nil {
		return file_domain.File{}, f.updateErr
	}
	f.updatedFolderID = folderID
	f.found.Filename = filename
	f.found.FolderID = folderID
	return f.found, nil
}

func (f *mockFileRepository) UpdateStatus(ctx context.Context, id uuid.UUID, userID int64, status file_domain.FileStatus) (file_domain.File, error) {
	f.found.Status = status
	return f.found, nil
}

func (f *mockFileRepository) UpdatePreview(ctx context.Context, id uuid.UUID, userID int64, status file_domain.FileStatus, previewPath *string, previewMimetype *string) (file_domain.File, error) {
	f.found.Status = status
	f.found.PreviewPath = previewPath
	f.found.PreviewMime = previewMimetype
	return f.found, nil
}

func (f *mockFileRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if f.deleteErr != nil {
		return f.deleteErr
	}
	f.deletedID = id
	return nil
}

type mockFileStorage struct {
	saveErr       error
	deleteErr     error
	openErr       error
	savePath      string
	savedFilename string
	savedAtPath   string
	savedAtBytes  int
	deleted       []string
	openedPath    string
	openContent   string
	openBytes     []byte
}

func (f *mockFileStorage) Save(ctx context.Context, filename string, src io.Reader) (string, error) {
	if f.saveErr != nil {
		return "", f.saveErr
	}
	f.savedFilename = filename
	if f.savePath == "" {
		return "objects/" + filename, nil
	}
	return f.savePath, nil
}

func (f *mockFileStorage) SaveAtPath(ctx context.Context, path string, src io.Reader) error {
	f.savedAtPath = path
	data, err := io.ReadAll(src)
	if err != nil {
		return err
	}
	f.savedAtBytes = len(data)
	return nil
}

func (f *mockFileStorage) Delete(ctx context.Context, path string) error {
	if f.deleteErr != nil {
		return f.deleteErr
	}
	f.deleted = append(f.deleted, path)
	return nil
}

func (f *mockFileStorage) Open(ctx context.Context, path string) (io.ReadCloser, error) {
	if f.openErr != nil {
		return nil, f.openErr
	}
	f.openedPath = path
	if f.openBytes != nil {
		return io.NopCloser(bytes.NewReader(f.openBytes)), nil
	}
	return io.NopCloser(strings.NewReader(f.openContent)), nil
}

type mockFileFolderRepository struct {
	findCalls int
	folder    folder_domain.Folder
	err       error
}

func (f *mockFileFolderRepository) FindByIDAndUserID(ctx context.Context, id uuid.UUID, userID int64) (folder_domain.Folder, error) {
	f.findCalls++
	if f.err != nil {
		return folder_domain.Folder{}, f.err
	}
	return f.folder, nil
}

type mockFileEventPublisher struct {
	uploaded []file_domain.File
	err      error
}

func (f *mockFileEventPublisher) PublishUploaded(ctx context.Context, file file_domain.File) error {
	if f.err != nil {
		return f.err
	}
	f.uploaded = append(f.uploaded, file)
	return nil
}

type mockFileMetrics struct {
	uploadedBytes int64
	started       int
	finished      []FileProcessingResult
}

func (f *mockFileMetrics) UploadSucceeded(sizeBytes int64) {
	f.uploadedBytes += sizeBytes
}

func (f *mockFileMetrics) ProcessingStarted() {
	f.started++
}

func (f *mockFileMetrics) ObserveProcessingFinished(result FileProcessingResult, duration time.Duration) {
	f.finished = append(f.finished, result)
}

func testJPEG(t *testing.T, width int, height int) []byte {
	t.Helper()

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{R: uint8(x % 255), G: uint8(y % 255), B: 120, A: 255})
		}
	}

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 80}); err != nil {
		t.Fatalf("encode jpeg: %v", err)
	}
	return buf.Bytes()
}
