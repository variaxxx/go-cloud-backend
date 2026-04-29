package folder_app

import (
	folder_domain "cloud/internal/features/folder/domain"
	"context"
)

type FolderUseCase interface {
	Create(
		ctx context.Context,
		name string,
		userID int64,
		parentID *int64,
	) (folder_domain.Folder, error)
}
