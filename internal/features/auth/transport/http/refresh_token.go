package auth_http

import (
	core_errors "cloud/internal/core/errors"
	core_logger "cloud/internal/core/logger"
	core_http_response "cloud/internal/core/transport/http/response"
	"errors"
	"fmt"
	"net/http"
)

type RefreshTokensResponse struct {
	AccessToken string `json:"access_token"`
}

func (h *Handler) RefreshTokens(rw http.ResponseWriter, r *http.Request) {
	log := core_logger.FromContext(r.Context())
	rh := core_http_response.NewHTTPResponseHandler(log, rw)

	refreshTokenCookie, err := r.Cookie("refresh_token")
	if err != nil {
		rh.ErrorResponse(
			fmt.Errorf("%w: refresh cookie missing", core_errors.ErrInvalidArgument),
			"Refresh token missing",
		)
		return
	}

	tokens, err := h.authService.RefreshTokens(r.Context(), refreshTokenCookie.Value)
	if err != nil {
		switch {
		case errors.Is(err, core_errors.ErrInvalidArgument), errors.Is(err, core_errors.ErrNotFound):
			rh.ErrorResponse(core_errors.ErrInvalidArgument, "Invalid refresh token")
		default:
			rh.ErrorResponse(err, "Unable to refresh tokens right now")
		}
		return
	}

	h.setRefreshTokenCookie(rw, tokens.RefreshToken)

	rh.JSONResponse(RefreshTokensResponse{
		AccessToken: tokens.AccessToken,
	}, http.StatusOK)
}
