package core_http_request

import (
	core_errors "cloud/internal/core/errors"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func DecodeAndValidateRequest(
	r *http.Request,
	dest any,
) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dest); err != nil {
		return fmt.Errorf("%w: decode JSON: %w", core_errors.ErrInvalidArgument, err)
	}

	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return fmt.Errorf("%w: request body must contain a single JSON object", core_errors.ErrInvalidArgument)
	}

	if err := validate.Struct(dest); err != nil {
		return fmt.Errorf("%w: validation: %w", core_errors.ErrInvalidArgument, err)
	}

	return nil
}
