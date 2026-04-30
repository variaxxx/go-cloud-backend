package folder_http

import (
	core_errors "cloud/internal/core/errors"
	core_logger "cloud/internal/core/logger"
	core_http_auth "cloud/internal/core/transport/http/auth"
	core_http_request "cloud/internal/core/transport/http/request"
	core_http_response "cloud/internal/core/transport/http/response"
	core_http_utils "cloud/internal/core/transport/http/utils"
	folder_app "cloud/internal/features/folder/application"
	"encoding/json"
	"errors"
	"net/http"
)

type NullableInt64Field struct {
	Present bool
	Value   *int64
}

func (f *NullableInt64Field) UnmarshalJSON(data []byte) error {
	f.Present = true

	if string(data) == "null" {
		f.Value = nil
		return nil
	}

	var value int64
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	f.Value = &value
	return nil
}

type EditRequest struct {
	Name     *string            `json:"name" validate:"omitempty,min=1,max=255"`
	ParentID NullableInt64Field `json:"parent_id"`
}

type EditResponse = FolderDTO

func (r EditRequest) ToEditFolderParams(
	folderID int64,
	userID int64,
) (folder_app.EditFolderParams, error) {
	params := folder_app.EditFolderParams{
		ID:     folderID,
		UserID: userID,
		Name:   r.Name,
	}

	if !r.ParentID.Present {
		return params, nil
	}

	params.IsParentIDUpdate = true

	if r.ParentID.Value != nil && *r.ParentID.Value <= 0 {
		return folder_app.EditFolderParams{}, core_errors.ErrInvalidArgument
	}

	params.ParentID = r.ParentID.Value

	return params, nil
}

func (h *Handler) Edit(rw http.ResponseWriter, r *http.Request) {
	log := core_logger.FromContext(r.Context())
	rh := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, ok := core_http_auth.UserIDFromContext(r.Context())
	if !ok {
		rh.ErrorResponse(core_errors.ErrUnauthorized, "Unauthorized")
		return
	}

	folderID, err := core_http_utils.GetIntPathValue(r, "id")
	if err != nil {
		rh.ErrorResponse(err, "Invalid folder ID")
		return
	}

	var request EditRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		rh.ErrorResponse(err, "Invalid request body")
		return
	}

	params, err := request.ToEditFolderParams(folderID, userID)
	if err != nil {
		rh.ErrorResponse(err, "Invalid request body")
		return
	}

	updatedFolder, err := h.folderService.Edit(
		r.Context(),
		params,
	)
	if err != nil {
		switch {
		case errors.Is(err, core_errors.ErrConflict):
			rh.ErrorResponse(err, "A folder with this name already exists in the selected directory")
			return
		case errors.Is(err, core_errors.ErrNotFound):
			rh.ErrorResponse(err, "Folder or parent folder not found")
			return
		case errors.Is(err, core_errors.ErrInvalidArgument):
			rh.ErrorResponse(err, "Invalid folder update data")
			return
		}

		rh.ErrorResponse(err, "Unable to edit folder")
		return
	}

	rh.JSONResponse(EditResponse(NewFolderDTO(updatedFolder)), http.StatusOK)
}
