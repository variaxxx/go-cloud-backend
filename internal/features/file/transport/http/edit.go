package file_http

import (
	core_errors "cloud/internal/core/errors"
	core_logger "cloud/internal/core/logger"
	core_http_auth "cloud/internal/core/transport/http/auth"
	core_http_request "cloud/internal/core/transport/http/request"
	core_http_response "cloud/internal/core/transport/http/response"
	core_http_utils "cloud/internal/core/transport/http/utils"
	file_app "cloud/internal/features/file/application"
	"errors"
	"net/http"

	"github.com/google/uuid"
)

type EditRequest struct {
	Filename *string                             `json:"filename" validate:"omitempty,min=1,max=255"`
	FolderID core_http_request.NullableUUIDField `json:"folder_id"`
}

type EditResponse = fileDTO

func (r EditRequest) ToEditFileParams(
	fileID uuid.UUID,
	userID int64,
) file_app.EditFileParams {
	params := file_app.EditFileParams{
		ID:       fileID,
		UserID:   userID,
		Filename: r.Filename,
	}

	if !r.FolderID.Present {
		return params
	}

	params.IsFolderIDUpdate = true
	params.FolderID = r.FolderID.Value

	return params
}

func (h *handler) Edit(rw http.ResponseWriter, r *http.Request) {
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

	var request EditRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		rh.ErrorResponse(err, "Invalid request body")
		return
	}

	updatedFile, err := h.fileService.Edit(
		r.Context(),
		request.ToEditFileParams(fileID, userID),
	)
	if err != nil {
		switch {
		case errors.Is(err, core_errors.ErrConflict):
			rh.ErrorResponse(err, "A file with this name already exists in the selected directory")
			return
		case errors.Is(err, core_errors.ErrNotFound):
			rh.ErrorResponse(err, "File or folder not found")
			return
		case errors.Is(err, core_errors.ErrInvalidArgument):
			rh.ErrorResponse(err, "Invalid file update data")
			return
		}

		rh.ErrorResponse(err, "Unable to edit file")
		return
	}

	rh.JSONResponse(EditResponse(NewFileDTO(updatedFile)), http.StatusOK)
}
