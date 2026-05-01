package user_http

import (
	core_errors "cloud/internal/core/errors"
	core_logger "cloud/internal/core/logger"
	core_http_auth "cloud/internal/core/transport/http/auth"
	core_http_response "cloud/internal/core/transport/http/response"
	"errors"
	"net/http"
)

type GetMeResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

func (h *handler) GetMe(rw http.ResponseWriter, r *http.Request) {
	log := core_logger.FromContext(r.Context())
	rh := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, ok := core_http_auth.UserIDFromContext(r.Context())
	if !ok {
		rh.ErrorResponse(core_errors.ErrUnauthorized, "Unauthorized")
		return
	}

	user, err := h.userService.GetMe(r.Context(), userID)
	if err != nil {
		if errors.Is(err, core_errors.ErrNotFound) {
			rh.ErrorResponse(err, "User not found")
			return
		}

		rh.ErrorResponse(err, "Unable to get profile info")
		return
	}

	rh.JSONResponse(GetMeResponse{
		ID:       user.ID,
		Username: user.Username,
	}, http.StatusOK)
}
