package core_http_middleware

import (
	"net/http"
	"strings"
)

func CORS(frontendURL string) Middleware {
	allowedOrigin := strings.TrimRight(frontendURL, "/")

	return func(next http.Handler) http.Handler {
		if allowedOrigin == "" {
			return next
		}

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			originAllowed := origin == allowedOrigin
			if originAllowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-Request-ID")
				w.Header().Add("Vary", "Origin")
			}

			if originAllowed && r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
