package folder_http

import (
	file_domain "cloud/internal/features/file/domain"
	folder_domain "cloud/internal/features/folder/domain"
	"time"

	"github.com/google/uuid"
)

type FolderDTO struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	ParentID  *uuid.UUID `json:"parent_id,omitempty"`
	CreatedAt string     `json:"created_at"`
	UpdatedAt string     `json:"updated_at"`
}

func NewFolderDTO(folder folder_domain.Folder) FolderDTO {
	return FolderDTO{
		ID:        folder.ID,
		Name:      folder.Name,
		ParentID:  folder.ParentID,
		CreatedAt: folder.CreatedAt.Format(time.RFC3339),
		UpdatedAt: folder.UpdatedAt.Format(time.RFC3339),
	}
}

func NewFolderDTOs(folders []folder_domain.Folder) []FolderDTO {
	dtos := make([]FolderDTO, 0, len(folders))
	for _, folder := range folders {
		dtos = append(dtos, NewFolderDTO(folder))
	}

	return dtos
}

type FileDTOView struct {
	ID        uuid.UUID `json:"id"`
	Filename  string    `json:"filename"`
	Mimetype  *string   `json:"mimetype,omitempty"`
	Status    string    `json:"status"`
	SizeBytes int64     `json:"size_bytes"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

func NewFileDTOView(file file_domain.File) FileDTOView {
	return FileDTOView{
		ID:        file.ID,
		Filename:  file.Filename,
		Mimetype:  file.Mimetype,
		Status:    string(file.Status),
		SizeBytes: file.SizeBytes,
		CreatedAt: file.CreatedAt.Format(time.RFC3339),
		UpdatedAt: file.UpdatedAt.Format(time.RFC3339),
	}
}

func NewFileDTOViews(files []file_domain.File) []FileDTOView {
	dtos := make([]FileDTOView, 0, len(files))
	for _, file := range files {
		dtos = append(dtos, NewFileDTOView(file))
	}

	return dtos
}
