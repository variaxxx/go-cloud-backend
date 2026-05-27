package auth_http

import "net/http"

func (h *handler) Logout(rw http.ResponseWriter, _ *http.Request) {
	h.clearRefreshTokenCookie(rw)
	rw.WriteHeader(http.StatusNoContent)
}
