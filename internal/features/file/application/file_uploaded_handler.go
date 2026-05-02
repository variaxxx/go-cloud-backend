package file_app

import (
	file_domain "cloud/internal/features/file/domain"
	"context"
	"fmt"
)

type FileUploadedEventHandler interface {
	HandleUploaded(
		ctx context.Context,
		event FileUploadedEvent,
	) error
}

type fileUploadedEventHandler struct {
	fileRepo file_domain.FileRepository
	storage  file_domain.FileStorage
}

func NewFileUploadedEventHandler(
	fileRepo file_domain.FileRepository,
	storage file_domain.FileStorage,
) *fileUploadedEventHandler {
	return &fileUploadedEventHandler{
		fileRepo: fileRepo,
		storage:  storage,
	}
}

func (h *fileUploadedEventHandler) HandleUploaded(
	ctx context.Context,
	event FileUploadedEvent,
) error {
	fmt.Println(event.Filename)

	return nil
}
