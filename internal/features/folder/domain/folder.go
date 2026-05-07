package folder_domain

import (
	"time"

	"github.com/google/uuid"
)

type Folder struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	UserID    int64
	ParentID  *uuid.UUID
}
