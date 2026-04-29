package folder_domain

import "context"

type FolderRepository interface {
	Create(
		ctx context.Context,
		name string,
		userID int64,
		parentID *int64,
	) (Folder, error)

	FindByIDAndUserID(
		ctx context.Context,
		id int64,
		userID int64,
	) (Folder, error)
}
