package file_http

import (
	file_domain "cloud/internal/features/file/domain"
	"time"

	"github.com/google/uuid"
)

type fileDTO struct {
	ID        uuid.UUID  `json:"id"`
	Filename  string     `json:"filename"`
	Mimetype  *string    `json:"mimetype,omitempty"`
	Status    string     `json:"status"`
	SizeBytes int64      `json:"size_bytes"`
	FolderID  *uuid.UUID `json:"folder_id,omitempty"`
	CreatedAt string     `json:"created_at"`
	UpdatedAt string     `json:"updated_at"`
}

func NewFileDTO(file file_domain.File) fileDTO {
	return fileDTO{
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

func NewFileDTOs(files []file_domain.File) []fileDTO {
	dtos := make([]fileDTO, 0, len(files))
	for _, file := range files {
		dtos = append(dtos, NewFileDTO(file))
	}

	return dtos
}
