package file_http

import (
	file_domain "cloud/internal/features/file/domain"
	"time"
)

type FileDTO struct {
	ID        int64   `json:"id"`
	Filename  string  `json:"filename"`
	Mimetype  *string `json:"mimetype,omitempty"`
	Status    string  `json:"status"`
	SizeBytes int64   `json:"size_bytes"`
	FolderID  *int64  `json:"folder_id,omitempty"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

func NewFileDTO(file file_domain.File) FileDTO {
	return FileDTO{
		ID:        file.ID,
		Filename:  file.Filename,
		Mimetype:  file.Mimetype,
		Status:    string(file.Status),
		SizeBytes: file.SizeBytes,
		FolderID:  file.FolderID,
		CreatedAt: file.CreatedAt.Format(time.RFC3339),
		UpdatedAt: file.UpdatedAt.Format(time.RFC3339),
	}
}

func NewFileDTOs(files []file_domain.File) []FileDTO {
	dtos := make([]FileDTO, 0, len(files))
	for _, file := range files {
		dtos = append(dtos, NewFileDTO(file))
	}

	return dtos
}
