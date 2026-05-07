package core_http_middleware

import (
	"net/http"

	"github.com/google/uuid"
)

const (
	requestIdHeader = "X-Request-ID"
)

func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIdHeader)
			if requestID == "" {
				requestID = uuid.NewString()
			}

			r.Header.Set(requestIdHeader, requestID)
			w.Header().Set(requestIdHeader, requestID)

			next.ServeHTTP(w, r)
		})
	}
}
