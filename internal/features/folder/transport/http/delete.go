package folder_http

import (
	core_errors "cloud/internal/core/errors"
	core_logger "cloud/internal/core/logger"
	core_http_response "cloud/internal/core/transport/http/response"
	core_http_utils "cloud/internal/core/transport/http/utils"
	auth_http "cloud/internal/features/auth/transport/http"
	"errors"
	"net/http"
)

func (h *Handler) Delete(rw http.ResponseWriter, r *http.Request) {
	log := core_logger.FromContext(r.Context())
	rh := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, ok := auth_http.UserIDFromContext(r.Context())
	if !ok {
		rh.ErrorResponse(core_errors.ErrUnauthorized, "Unauthorized")
		return
	}

	folderID, err := core_http_utils.GetIntPathValue(r, "id")
	if err != nil {
		rh.ErrorResponse(err, "Invalid folder ID")
		return
	}

	if err := h.folderService.Delete(r.Context(), folderID, userID); err != nil {
		switch {
		case errors.Is(err, core_errors.ErrNotFound):
			rh.ErrorResponse(err, "Folder not found")
			return
		}

		rh.ErrorResponse(err, "Unable to delete folder")
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
