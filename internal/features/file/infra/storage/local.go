package file_storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type localFileStorage struct {
	baseDir string
}

func NewLocalFileStorage(
	baseDir string,
) *localFileStorage {
	return &localFileStorage{
		baseDir: baseDir,
	}
}

func (s *localFileStorage) Save(
	ctx context.Context,
	filename string,
	src io.Reader,
) (path string, err error) {
	ext := filepath.Ext(filename)
	name := uuid.NewString() + ext

	relativePath := filepath.Join(name)
	fullPath := filepath.Join(s.baseDir, relativePath)

	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return "", fmt.Errorf("create storage dir: %w", err)
	}

	dst, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("create dest file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("copy file data: %w", err)
	}

	return relativePath, nil
}

func (s *localFileStorage) SaveAtPath(
	ctx context.Context,
	path string,
	src io.Reader,
) error {
	fullPath := filepath.Join(s.baseDir, path)

	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return fmt.Errorf("create storage dir: %w", err)
	}

	dst, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("create dest file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("copy file data: %w", err)
	}

	return nil
}

func (s *localFileStorage) Delete(
	ctx context.Context,
	path string,
) error {
	fullPath := filepath.Join(s.baseDir, path)
	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("delete file from storage: %w", err)
	}

	return nil
}

func (s *localFileStorage) Open(
	ctx context.Context,
	path string,
) (io.ReadCloser, error) {
	fullPath := filepath.Join(s.baseDir, path)

	file, err := os.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("open file from storage: %w", err)
	}

	return file, nil
}
