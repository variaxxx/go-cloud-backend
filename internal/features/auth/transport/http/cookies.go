package auth_http

import "net/http"

func (h *Handler) setRefreshTokenCookie(
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
