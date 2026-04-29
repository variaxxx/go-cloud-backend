package file_domain

import "context"

type FileRepository interface {
	Create(
		ctx context.Context,
		filename string,
		mimetype *string,
		status FileStatus,
		storagePath string,
		sizeBytes int64,
		userID int64,
	) (File, error)
}
