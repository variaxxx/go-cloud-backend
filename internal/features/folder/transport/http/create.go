package folder_http

import (
	core_errors "cloud/internal/core/errors"
	core_logger "cloud/internal/core/logger"
	core_http_auth "cloud/internal/core/transport/http/auth"
	core_http_request "cloud/internal/core/transport/http/request"
	core_http_response "cloud/internal/core/transport/http/response"
	folder_app "cloud/internal/features/folder/application"
	"errors"
	"net/http"

	"github.com/google/uuid"
)

type CreateRequest struct {
	Name     string     `json:"name" validate:"required,min=1,max=255"`
	ParentID *uuid.UUID `json:"parent_id"`
}

type CreateResponse = FolderDTO

func (h *Handler) Create(rw http.ResponseWriter, r *http.Request) {
	log := core_logger.FromContext(r.Context())
	rh := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, ok := core_http_auth.UserIDFromContext(r.Context())
	if !ok {
		rh.ErrorResponse(core_errors.ErrUnauthorized, "Unauthorized")
		return
	}

	var request CreateRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		rh.ErrorResponse(err, "Invalid request body")
		return
	}

	createdFolder, err := h.folderService.Create(
		r.Context(),
		folder_app.CreateFolderParams{
			UserID:   userID,
			Name:     request.Name,
			ParentID: request.ParentID,
		},
	)
	if err != nil {
		switch {
		case errors.Is(err, core_errors.ErrNotFound):
			rh.ErrorResponse(err, "Parent folder not found")
			return
		case errors.Is(err, core_errors.ErrConflict):
			rh.ErrorResponse(err, "A folder with this name already exists in the selected directory")
			return
		}

		rh.ErrorResponse(err, "Unable to create folder")
		return
	}

	rh.JSONResponse(CreateResponse(NewFolderDTO(createdFolder)), http.StatusCreated)
}
