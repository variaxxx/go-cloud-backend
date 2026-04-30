package core_http_utils

import (
	core_errors "cloud/internal/core/errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func GetUUIDPathValue(r *http.Request, key string) (uuid.UUID, error) {
	pathValue := r.PathValue(key)
	if pathValue == "" {
		return uuid.Nil, fmt.Errorf(
			"%w: no key='%s' in path values",
			core_errors.ErrInvalidArgument,
			key,
		)
	}

	val, err := uuid.Parse(pathValue)
	if err != nil {
		return uuid.Nil, fmt.Errorf(
			"%w: path value='%s' by key='%s' is not a valid UUID: %v",
			core_errors.ErrInvalidArgument, pathValue, key, err,
		)
	}

	return val, nil
}
