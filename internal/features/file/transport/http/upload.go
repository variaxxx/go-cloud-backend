package file_http

import (
	core_errors "cloud/internal/core/errors"
	core_logger "cloud/internal/core/logger"
	core_http_auth "cloud/internal/core/transport/http/auth"
	core_http_response "cloud/internal/core/transport/http/response"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

type UploadResponse = FileDTO

func (h *Handler) Upload(rw http.ResponseWriter, r *http.Request) {
	log := core_logger.FromContext(r.Context())
	rh := core_http_response.NewHTTPResponseHandler(log, rw)

	// Size limit to 10MB
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		rh.ErrorResponse(
			fmt.Errorf("%w: parse multipart form: %w", core_errors.ErrInvalidArgument, err),
			"Request body must be multipart/form-data with a file field",
		)
		return
	}

	userID, ok := core_http_auth.UserIDFromContext(r.Context())
	if !ok {
		rh.ErrorResponse(core_errors.ErrUnauthorized, "Unauthorized")
		return
	}

	var folderID *int64
	if value := r.FormValue("folder_id"); value != "" {
		parsedFolderID, err := strconv.ParseInt(value, 10, 64)
		if err != nil || parsedFolderID <= 0 {
			rh.ErrorResponse(core_errors.ErrInvalidArgument, "folder_id must be a positive integer")
			return
		}

		folderID = &parsedFolderID
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		rh.ErrorResponse(core_errors.ErrInvalidArgument, "File field is missing")
		return
	}
	defer file.Close()

	var mimetype *string
	if value := header.Header.Get("Content-Type"); value != "" {
		mimetype = &value
	}

	uploadedFile, err := h.fileService.Upload(
		r.Context(),
		header.Filename,
		mimetype,
		header.Size,
		userID,
		folderID,
		file,
	)
	if err != nil {
		switch {
		case errors.Is(err, core_errors.ErrNotFound):
			rh.ErrorResponse(err, "Folder not found")
			return
		case errors.Is(err, core_errors.ErrConflict):
			rh.ErrorResponse(err, "A file with this name already exists in the selected folder")
			return
		}

		rh.ErrorResponse(err, "Unable to upload file")
		return
	}

	rh.JSONResponse(UploadResponse(NewFileDTO(uploadedFile)), http.StatusCreated)
}
