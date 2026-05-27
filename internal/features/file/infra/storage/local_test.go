package file_storage

import (
	"context"
	"io"
	"strings"
	"testing"
)

func TestLocalFileStorageSaveOpenSaveAtPathAndDelete(t *testing.T) {
	ctx := context.Background()
	storage := NewLocalFileStorage(t.TempDir())

	path, err := storage.Save(ctx, "document.txt", strings.NewReader("hello"))
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}
	if !strings.HasSuffix(path, ".txt") {
		t.Fatalf("Save() path = %q, want .txt suffix", path)
	}

	content, err := storage.Open(ctx, path)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	body, err := io.ReadAll(content)
	content.Close()
	if err != nil {
		t.Fatalf("ReadAll() error = %v", err)
	}
	if string(body) != "hello" {
		t.Fatalf("Open() body = %q, want hello", string(body))
	}

	if err := storage.SaveAtPath(ctx, "previews/document.jpg", strings.NewReader("preview")); err != nil {
		t.Fatalf("SaveAtPath() error = %v", err)
	}
	preview, err := storage.Open(ctx, "previews/document.jpg")
	if err != nil {
		t.Fatalf("Open(preview) error = %v", err)
	}
	previewBody, err := io.ReadAll(preview)
	preview.Close()
	if err != nil {
		t.Fatalf("ReadAll(preview) error = %v", err)
	}
	if string(previewBody) != "preview" {
		t.Fatalf("preview body = %q, want preview", string(previewBody))
	}

	if err := storage.Delete(ctx, path); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	if err := storage.Delete(ctx, path); err != nil {
		t.Fatalf("Delete() missing file error = %v, want nil", err)
	}
}
