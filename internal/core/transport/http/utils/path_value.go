package core_http_utils

import (
	core_errors "cloud/internal/core/errors"
	"fmt"
	"net/http"
	"strconv"
)

func GetIntPathValue(r *http.Request, key string) (int64, error) {
	pathValue := r.PathValue(key)
	if pathValue == "" {
		return 0, fmt.Errorf(
			"%w: no key='%s' in path values",
			core_errors.ErrInvalidArgument,
			key,
		)
	}

	val, err := strconv.ParseInt(pathValue, 10, 64)
	if err != nil {
		return 0, fmt.Errorf(
			"%w: path value='%s' by key='%s' is not a valid integer: %v",
			core_errors.ErrInvalidArgument, pathValue, key, err,
		)
	}

	return val, nil
}
