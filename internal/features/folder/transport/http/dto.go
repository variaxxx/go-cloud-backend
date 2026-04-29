package folder_http

import (
	folder_domain "cloud/internal/features/folder/domain"
	"time"
)

type FolderDTO struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	ParentID  *int64 `json:"parent_id,omitempty"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
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
