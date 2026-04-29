package file_domain

import (
	"context"
	"io"
)

type FileStorage interface {
	Save(
		ctx context.Context,
		filename string,
		src io.Reader,
	) (path string, err error)

	Delete(
		ctx context.Context,
		path string,
	) error
}
