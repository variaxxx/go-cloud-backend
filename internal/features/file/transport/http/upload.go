package file_http

import (
	core_errors "cloud/internal/core/errors"
	core_logger "cloud/internal/core/logger"
	core_http_response "cloud/internal/core/transport/http/response"
	auth_http "cloud/internal/features/auth/transport/http"
	"fmt"
	"net/http"
)

type UploadResponse struct {
	ID          int64   `json:"id"`
	Filename    string  `json:"filename"`
	Mimetype    *string `json:"mimetype,omitempty"`
	Status      string  `json:"status"`
	StoragePath string  `json:"storage_path"`
	SizeBytes   int64   `json:"size_bytes"`
	UserID      int64   `json:"user_id"`
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
		file,
	)
	if err != nil {
		rh.ErrorResponse(err, "Unable to upload file")
		return
	}

	rh.JSONResponse(UploadResponse{
		ID:          uploadedFile.ID,
		Filename:    uploadedFile.Filename,
		Mimetype:    uploadedFile.Mimetype,
		Status:      string(uploadedFile.Status),
		StoragePath: uploadedFile.StoragePath,
		SizeBytes:   uploadedFile.SizeBytes,
		UserID:      uploadedFile.UserID,
	}, http.StatusCreated)
}
