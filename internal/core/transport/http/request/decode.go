package core_http_request

import (
	core_errors "cloud/internal/core/errors"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func DecodeAndValidateRequest(
	r *http.Request,
	dest any,
) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return fmt.Errorf("%w: decode JSON: %w", core_errors.ErrInvalidArgument, err)
	}

	if err := validate.Struct(dest); err != nil {
		return fmt.Errorf("%w: validation: %w", core_errors.ErrInvalidArgument, err)
	}

	return nil
}
