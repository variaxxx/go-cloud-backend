package folder_http

import (
	file_domain "cloud/internal/features/file/domain"
	folder_domain "cloud/internal/features/folder/domain"
	"time"

	"github.com/google/uuid"
)

type folderDTO struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	ParentID  *uuid.UUID `json:"parent_id,omitempty"`
	CreatedAt string     `json:"created_at"`
	UpdatedAt string     `json:"updated_at"`
}

type folderPathDTO struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func NewFolderDTO(folder folder_domain.Folder) folderDTO {
	return folderDTO{
		ID:        folder.ID,
		Name:      folder.Name,
		ParentID:  folder.ParentID,
		CreatedAt: folder.CreatedAt.Format(time.RFC3339),
		UpdatedAt: folder.UpdatedAt.Format(time.RFC3339),
	}
}

func NewFolderDTOs(folders []folder_domain.Folder) []folderDTO {
	dtos := make([]folderDTO, 0, len(folders))
	for _, folder := range folders {
		dtos = append(dtos, NewFolderDTO(folder))
	}

	return dtos
}

func NewFolderPathDTO(folder folder_domain.Folder) folderPathDTO {
	return folderPathDTO{
		ID:   folder.ID,
		Name: folder.Name,
	}
}

func NewFolderPathDTOs(folders []folder_domain.Folder) []folderPathDTO {
	dtos := make([]folderPathDTO, 0, len(folders))
	for _, folder := range folders {
		dtos = append(dtos, NewFolderPathDTO(folder))
	}

	return dtos
}

type fileDTOView struct {
	ID         uuid.UUID `json:"id"`
	Filename   string    `json:"filename"`
	Mimetype   *string   `json:"mimetype,omitempty"`
	Status     string    `json:"status"`
	SizeBytes  int64     `json:"size_bytes"`
	HasPreview bool      `json:"has_preview"`
	CreatedAt  string    `json:"created_at"`
	UpdatedAt  string    `json:"updated_at"`
}

func NewFileDTOView(file file_domain.File) fileDTOView {
	return fileDTOView{
		ID:         file.ID,
		Filename:   file.Filename,
		Mimetype:   file.Mimetype,
		Status:     string(file.Status),
		SizeBytes:  file.SizeBytes,
		HasPreview: file.PreviewPath != nil,
		CreatedAt:  file.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  file.UpdatedAt.Format(time.RFC3339),
	}
}

func NewFileDTOViews(files []file_domain.File) []fileDTOView {
	dtos := make([]fileDTOView, 0, len(files))
	for _, file := range files {
		dtos = append(dtos, NewFileDTOView(file))
	}

	return dtos
}
