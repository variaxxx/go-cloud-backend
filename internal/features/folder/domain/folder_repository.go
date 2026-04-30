package folder_domain

import "context"

type FolderRepository interface {
	Create(
		ctx context.Context,
		name string,
		userID int64,
		parentID *int64,
	) (Folder, error)

	Update(
		ctx context.Context,
		id int64,
		userID int64,
		name string,
		parentID *int64,
	) (Folder, error)

	FindByIDAndUserID(
		ctx context.Context,
		id int64,
		userID int64,
	) (Folder, error)

	FindByIDAndUserIDWithPath(
		ctx context.Context,
		id int64,
		userID int64,
	) (Folder, []string, error)

	FindByParentIDAndUserID(
		ctx context.Context,
		parentID *int64,
		userID int64,
	) ([]Folder, error)

	WouldCreateCycle(
		ctx context.Context,
		id int64,
		parentID int64,
		userID int64,
	) (bool, error)

	Delete(
		ctx context.Context,
		id int64,
	) error
}
