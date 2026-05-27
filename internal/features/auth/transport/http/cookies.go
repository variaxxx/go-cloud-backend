package auth_http

import (
	"net/http"
	"time"
)

func (h *handler) setRefreshTokenCookie(
	rw http.ResponseWriter,
	refreshToken string,
) {
	http.SetCookie(rw, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/auth",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(h.refreshTTL.Seconds()),
	})
}

func (h *handler) clearRefreshTokenCookie(
	rw http.ResponseWriter,
) {
	http.SetCookie(rw, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/auth",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
	})
}
