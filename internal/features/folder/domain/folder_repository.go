package folder_domain

import (
	"context"

	"github.com/google/uuid"
)

type FolderRepository interface {
	Create(
		ctx context.Context,
		name string,
		userID int64,
		parentID *uuid.UUID,
	) (Folder, error)

	Update(
		ctx context.Context,
		id uuid.UUID,
		userID int64,
		name string,
		parentID *uuid.UUID,
	) (Folder, error)

	FindByIDAndUserID(
		ctx context.Context,
		id uuid.UUID,
		userID int64,
	) (Folder, error)

	FindByIDAndUserIDWithPath(
		ctx context.Context,
		id uuid.UUID,
		userID int64,
	) (Folder, []string, error)

	FindByParentIDAndUserID(
		ctx context.Context,
		parentID *uuid.UUID,
		userID int64,
	) ([]Folder, error)

	WouldCreateCycle(
		ctx context.Context,
		id uuid.UUID,
		parentID uuid.UUID,
		userID int64,
	) (bool, error)

	Delete(
		ctx context.Context,
		id uuid.UUID,
	) error
}
