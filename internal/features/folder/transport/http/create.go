package folder_http

import (
	core_errors "cloud/internal/core/errors"
	core_logger "cloud/internal/core/logger"
	core_http_request "cloud/internal/core/transport/http/request"
	core_http_response "cloud/internal/core/transport/http/response"
	auth_http "cloud/internal/features/auth/transport/http"
	"errors"
	"net/http"
	"time"
)

type CreateRequest struct {
	Name     string `json:"name" validate:"required,min=1,max=255"`
	ParentID *int64 `json:"parent_id"`
}

type CreateResponse struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	ParentID  *int64 `json:"parent_id,omitempty"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (h *Handler) Create(rw http.ResponseWriter, r *http.Request) {
	log := core_logger.FromContext(r.Context())
	rh := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, ok := auth_http.UserIDFromContext(r.Context())
	if !ok {
		rh.ErrorResponse(core_errors.ErrUnauthorized, "Unauthorized")
		return
	}

	var request CreateRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		rh.ErrorResponse(err, "Invalid request body")
		return
	}

	createdFolder, err := h.folderService.Create(r.Context(), request.Name, userID, request.ParentID)
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

	rh.JSONResponse(CreateResponse{
		ID:        createdFolder.ID,
		Name:      createdFolder.Name,
		ParentID:  createdFolder.ParentID,
		CreatedAt: createdFolder.CreatedAt.Format(time.RFC3339),
		UpdatedAt: createdFolder.UpdatedAt.Format(time.RFC3339),
	}, http.StatusCreated)
}
