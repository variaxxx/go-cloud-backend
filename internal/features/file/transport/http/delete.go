package file_http

import (
	core_errors "cloud/internal/core/errors"
	core_logger "cloud/internal/core/logger"
	core_http_auth "cloud/internal/core/transport/http/auth"
	core_http_response "cloud/internal/core/transport/http/response"
	core_http_utils "cloud/internal/core/transport/http/utils"
	"errors"
	"net/http"
)

func (h *Handler) Delete(rw http.ResponseWriter, r *http.Request) {
	log := core_logger.FromContext(r.Context())
	rh := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, ok := core_http_auth.UserIDFromContext(r.Context())
	if !ok {
		rh.ErrorResponse(core_errors.ErrUnauthorized, "Unauthorized")
		return
	}

	fileID, err := core_http_utils.GetUUIDPathValue(r, "id")
	if err != nil {
		rh.ErrorResponse(err, "Invalid file ID")
		return
	}

	if err := h.fileService.Delete(r.Context(), fileID, userID); err != nil {
		switch {
		case errors.Is(err, core_errors.ErrNotFound):
			rh.ErrorResponse(err, "File not found")
			return
		}

		rh.ErrorResponse(err, "Unable to delete file")
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
