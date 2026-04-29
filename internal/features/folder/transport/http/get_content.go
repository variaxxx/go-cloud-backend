package folder_http

import (
	core_errors "cloud/internal/core/errors"
	core_logger "cloud/internal/core/logger"
	core_http_response "cloud/internal/core/transport/http/response"
	auth_http "cloud/internal/features/auth/transport/http"
	file_http "cloud/internal/features/file/transport/http"
	folder_app "cloud/internal/features/folder/application"
	"errors"
	"net/http"
	"strconv"
)

type GetContentResponse struct {
	FolderName *string             `json:"folder_name,omitempty"`
	FolderPath []string            `json:"folder_path"`
	Folders    []FolderDTO         `json:"folders"`
	Files      []file_http.FileDTO `json:"files"`
}

func (h *Handler) GetContent(rw http.ResponseWriter, r *http.Request) {
	log := core_logger.FromContext(r.Context())
	rh := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, ok := auth_http.UserIDFromContext(r.Context())
	if !ok {
		rh.ErrorResponse(core_errors.ErrUnauthorized, "Unauthorized")
		return
	}

	var folderID *int64
	if value := r.PathValue("id"); value != "" {
		parsedFolderID, err := strconv.ParseInt(value, 10, 64)
		if err != nil || parsedFolderID <= 0 {
			rh.ErrorResponse(core_errors.ErrInvalidArgument, "Folder id must be a positive integer")
			return
		}

		folderID = &parsedFolderID
	}

	content, err := h.folderService.GetContent(
		r.Context(),
		folder_app.GetFolderContentParams{
			UserID:   userID,
			FolderID: folderID,
		},
	)
	if err != nil {
		switch {
		case errors.Is(err, core_errors.ErrNotFound):
			rh.ErrorResponse(err, "Folder not found")
		default:
			rh.ErrorResponse(err, "Unable to get folder content")
		}
		return
	}

	rh.JSONResponse(GetContentResponse{
		FolderName: content.FolderName,
		FolderPath: content.FolderPath,
		Folders:    NewFolderDTOs(content.Folders),
		Files:      file_http.NewFileDTOs(content.Files),
	}, http.StatusOK)
}
