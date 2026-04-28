package auth_http

import (
	core_errors "cloud/internal/core/errors"
	core_logger "cloud/internal/core/logger"
	core_http_request "cloud/internal/core/transport/http/request"
	core_http_response "cloud/internal/core/transport/http/response"
	"errors"
	"net/http"
)

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=100"`
	Password string `json:"password" validate:"required,min=6,max=100"`
}

type RegisterResponse struct {
	AccessToken string `json:"access_token"`
}

func (h *Handler) Register(rw http.ResponseWriter, r *http.Request) {
	log := core_logger.FromContext(r.Context())
	rh := core_http_response.NewHTTPResponseHandler(log, rw)

	var request RegisterRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		rh.ErrorResponse(err, "Validation failed")
		return
	}

	tokens, err := h.authService.Register(r.Context(), request.Username, request.Password)
	if err != nil {
		if errors.Is(err, core_errors.ErrConflict) {
			rh.ErrorResponse(err, "User already exists")
			return
		}

		rh.ErrorResponse(err, "Register failed")
		return
	}

	http.SetCookie(rw, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		Path:     "/auth",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(h.refreshTTL.Seconds()),
	})

	rh.JSONResponse(RegisterResponse{
		AccessToken: tokens.AccessToken,
	}, http.StatusCreated)
}
