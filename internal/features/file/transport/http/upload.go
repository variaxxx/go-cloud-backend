package file_http

import (
	core_errors "cloud/internal/core/errors"
	core_logger "cloud/internal/core/logger"
	core_http_response "cloud/internal/core/transport/http/response"
	auth_http "cloud/internal/features/auth/transport/http"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

type UploadResponse struct {
	ID        int64   `json:"id"`
	Filename  string  `json:"filename"`
	Mimetype  *string `json:"mimetype,omitempty"`
	Status    string  `json:"status"`
	SizeBytes int64   `json:"size_bytes"`
	FolderID  *int64  `json:"folder_id,omitempty"`
}

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

	userID, ok := auth_http.UserIDFromContext(r.Context())
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

	rh.JSONResponse(UploadResponse{
		ID:        uploadedFile.ID,
		Filename:  uploadedFile.Filename,
		Mimetype:  uploadedFile.Mimetype,
		Status:    string(uploadedFile.Status),
		SizeBytes: uploadedFile.SizeBytes,
		FolderID:  uploadedFile.FolderID,
	}, http.StatusCreated)
}
