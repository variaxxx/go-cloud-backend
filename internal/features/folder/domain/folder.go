package folder_domain

import "time"

type Folder struct {
	ID        int64
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	UserID    int64
	ParentID  *int64
}
