package folder_http

import (
	core_errors "cloud/internal/core/errors"
	core_logger "cloud/internal/core/logger"
	core_http_auth "cloud/internal/core/transport/http/auth"
	core_http_response "cloud/internal/core/transport/http/response"
	"errors"
	"net/http"

	"github.com/google/uuid"
)

type GetContentResponse struct {
	FolderName *string       `json:"folder_name,omitempty"`
	FolderPath []string      `json:"folder_path"`
	Folders    []FolderDTO   `json:"folders"`
	Files      []FileDTOView `json:"files"`
}

func (h *Handler) GetContent(rw http.ResponseWriter, r *http.Request) {
	log := core_logger.FromContext(r.Context())
	rh := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, ok := core_http_auth.UserIDFromContext(r.Context())
	if !ok {
		rh.ErrorResponse(core_errors.ErrUnauthorized, "Unauthorized")
		return
	}

	var folderID *uuid.UUID
	if value := r.PathValue("id"); value != "" {
		parsedFolderID, err := uuid.Parse(value)
		if err != nil {
			rh.ErrorResponse(core_errors.ErrInvalidArgument, "Folder id must be a valid UUID")
			return
		}

		folderID = &parsedFolderID
	}

	content, err := h.folderService.GetContent(
		r.Context(),
		userID,
		folderID,
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
		Files:      NewFileDTOViews(content.Files),
	}, http.StatusOK)
}
