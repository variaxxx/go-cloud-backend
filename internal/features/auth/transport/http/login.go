package auth_http

import (
	core_errors "cloud/internal/core/errors"
	core_logger "cloud/internal/core/logger"
	core_http_request "cloud/internal/core/transport/http/request"
	core_http_response "cloud/internal/core/transport/http/response"
	"errors"
	"net/http"
)

type LoginRequest struct {
	Username string `json:"username" validate:"required,min=3,max=100"`
	Password string `json:"password" validate:"required,min=6,max=100"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

func (h *handler) Login(rw http.ResponseWriter, r *http.Request) {
	log := core_logger.FromContext(r.Context())
	rh := core_http_response.NewHTTPResponseHandler(log, rw)

	var request LoginRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		rh.ErrorResponse(err, "Request body is invalid")
		return
	}

	tokens, err := h.authService.Login(r.Context(), request.Username, request.Password)
	if err != nil {
		switch {
		case errors.Is(err, core_errors.ErrNotFound), errors.Is(err, core_errors.ErrInvalidArgument):
			rh.ErrorResponse(core_errors.ErrUnauthorized, "Invalid username or password")
		default:
			rh.ErrorResponse(err, "Unable to sign in right now")
		}
		return
	}

	h.setRefreshTokenCookie(rw, tokens.RefreshToken)

	rh.JSONResponse(LoginResponse{
		AccessToken: tokens.AccessToken,
	}, http.StatusOK)
}
