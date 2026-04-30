package file_http

import (
	core_errors "cloud/internal/core/errors"
	core_logger "cloud/internal/core/logger"
	core_http_auth "cloud/internal/core/transport/http/auth"
	core_http_response "cloud/internal/core/transport/http/response"
	core_http_utils "cloud/internal/core/transport/http/utils"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"go.uber.org/zap"
)

func (h *Handler) Download(rw http.ResponseWriter, r *http.Request) {
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

	result, err := h.fileService.Download(r.Context(), fileID, userID)
	if err != nil {
		switch {
		case errors.Is(err, core_errors.ErrNotFound):
			rh.ErrorResponse(err, "File not found")
			return
		}

		rh.ErrorResponse(err, "Unable to download file")
		return
	}
	defer result.Content.Close()

	contentType := "application/octet-stream"
	if result.Mimetype != nil && *result.Mimetype != "" {
		contentType = *result.Mimetype
	}

	rw.Header().Set("Content-Type", contentType)
	rw.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", strconv.Quote(result.Filename)))
	rw.WriteHeader(http.StatusOK)

	if _, err := io.Copy(rw, result.Content); err != nil {
		log.Error("Write file download response", zap.Error(err))
	}
}
